package p2c

import (
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
	"google.golang.org/grpc/resolver"
	"math"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

const (
	Name = "p2c"
	CPUUsage = "cpu_usage"
)

var (
	penalty = uint64(1000 * time.Millisecond * 250)
	forceGap = int64(time.Second * 3)
	tau = int64(600 * time.Millisecond)
)

func init(){
	balancer.Register(newBalance())
}

func newBalance() balancer.Builder{
	return base.NewBalancerBuilder(Name , &p2cPickerBuilder{} , base.Config{HealthCheck: true})
}

type subConn struct {
	conn balancer.SubConn
	addr resolver.Address

	lag uint64 // 表示该连接的延迟
	success uint64 // 表示该连接的成功率
	inflight int64 // 表示该连接上 正在请求的连接数
	serCpu uint64 // 表示该 连接对应node上 cpu的使用率
	stamp int64 // 上一次该 连接完成请求的时间 也就是上次 pick 执行 doneInfo的时候
	pick int64  // 上次 被选中的时间
	reqs int64 // 该连接 上一共经历的请求数 包括成功的和未成功的
}

func (s *subConn) valid() bool {
	// 表示 成功率 大于一半 而且 node cpu 使用率小于 900
	return s.health() > 500 && atomic.LoadUint64(&s.serCpu) < 900
}

func (s *subConn) health()uint64{
	return atomic.LoadUint64(&s.success)
}

func (s *subConn) load() uint64 {
	// 这里的计算公式就是  serCPU * 根号下lag * inflight
	// 值越大 说明 负载越高
	lag := uint64(math.Sqrt(float64(atomic.LoadUint64(&s.lag)))+1)
	load := lag * atomic.LoadUint64(&s.serCpu) * uint64(atomic.LoadInt64(&s.inflight))

	if load == 0{
		// 这里给的是一个默认值
		load =penalty
	}
	return load
}

type p2cPicker struct{
	subConns []*subConn
	lk sync.Mutex
	r *rand.Rand
	logTs int64
}


func (p *p2cPicker) Pick(info balancer.PickInfo) (balancer.PickResult, error){
	subConn , f , err := p.pick(info)
	res := balancer.PickResult{
		SubConn: nil,
		Done:    nil,
	}
	if err != nil {
		return res  ,err
	}
	res.SubConn = subConn
	res.Done = f
	return res , nil
}

func (p *p2cPicker)prePick() (nodeA *subConn ,nodeB *subConn) {
	// 这里选择的 node 若总是不健康的 只会选择3次
	for i:= 0 ; i < 3 ; i++{
		p.lk.Lock()
		a := p.r.Intn(len(p.subConns))
		b := p.r.Intn(len(p.subConns) - 1)
		p.lk.Unlock()
		if b == a {
			b = b + 1
		}
		nodeA , nodeB = p.subConns[a] , p.subConns[b]
		if nodeA.valid() || nodeB.valid(){
			break
		}

	}
	return
}

func (p *p2cPicker) pick(info balancer.PickInfo) (balancer.SubConn ,func(doneInfo balancer.DoneInfo) ,error) {
	// 随机出来的两条连接 最终pc是选择的那条
	var pc , uPc *subConn
	start := time.Now().UnixNano()
	if len(p.subConns) == 0{
		return nil , nil , balancer.ErrNoSubConnAvailable
	}else if len(p.subConns) == 1 {
		pc = p.subConns[0]
	}else{
		nodeA , nodeB := p.prePick()
		// 比较两个的 负载值
		// load 的值越大 说明 负载越高
		if nodeA.load() * nodeB.health() > nodeB.load() * nodeA.health(){
			pc = nodeB
			uPc = nodeA
		}else{
			pc = nodeA
			uPc = nodeB
		}
		// 这里检查 uPc的值 若uPc超过3秒没有被选中 那么强制选中
		pick := atomic.LoadInt64(&uPc.pick)
		if start - pick > forceGap && atomic.CompareAndSwapInt64(&uPc.pick , uPc.pick , start) {
			pc = uPc
		}

		if pc != uPc{
			atomic.StoreInt64(&pc.pick , start)
		}
		atomic.AddInt64(&pc.inflight , 1)
		atomic.AddInt64(&pc.reqs , 1)
	}
	return pc.conn ,func(doneInfo balancer.DoneInfo){
		atomic.AddInt64(&pc.inflight , -1)

		now := time.Now().UnixNano()
		oldStamp := atomic.SwapInt64(&pc.stamp , now)
		// 根据 EWMA 预测 lag 和 success 公式为： Vt=βVt+(1−β)θt
		// 计算  β = exp( （-(now-stmp)) / tau )
		w := math.Exp( float64( (-(now-oldStamp)) / tau ))
		lag := now - start
		if lag < 0 {
			lag = 0
		}
		oldLag := atomic.LoadUint64(&pc.lag)
		if oldLag == 0{
			w = 0.0
		}

		lag =int64( float64(oldLag) * w + (1-w)*float64(lag) )
		// 这里类型转换的时候 都是将 小类型往上转
		atomic.StoreUint64(&pc.lag , uint64(lag))

		success := uint64(1000)
		if doneInfo.Err != nil {
			// 这里若是调用失败 则将success的值 变为0
			success = 0
		}
		oldSuccess := atomic.LoadUint64(&pc.success)
		success = uint64( float64(oldSuccess) * w + (1-w) * float64(success) )
		atomic.StoreUint64(&pc.success , success)

		if cpuStr , ok := doneInfo.Trailer[CPUUsage] ; ok {
			if cpu , err := strconv.ParseUint(cpuStr[0] , 10 , 64); err != nil && cpu > 0 {
				atomic.StoreUint64(&pc.serCpu , cpu)
			}
		}
		// 这里3s清空一次 reqs
		logTs := atomic.LoadInt64(&p.logTs)
		if now - logTs > int64(time.Second * 3){
			if atomic.CompareAndSwapInt64(&p.logTs , logTs , now) {
				atomic.StoreInt64(&pc.reqs , 0)
			}
		}
	} , nil
}

type p2cPickerBuilder struct {}

func (p * p2cPickerBuilder) Build(info base.PickerBuildInfo) balancer.Picker{
	picker := &p2cPicker{
		subConns: make([]*subConn , 0 , len(info.ReadySCs)),
		lk:       sync.Mutex{},
		r:        rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for sc , addr := range info.ReadySCs{
		s := sc
		subC := &subConn{
			conn:     s,
			addr:     addr.Address,
			lag:      0,
			success:  100,
			inflight: 1,
			serCpu:   500,
		}
		picker.subConns = append(picker.subConns , subC)
	}
	return picker
}


