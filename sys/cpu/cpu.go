package cpu

import (
	"sync/atomic"
	"time"
)

const (
	interval = time.Millisecond * 500
)

var (
	usage uint64
	stats CPU
)

type CPU interface {
	Usage() (uint64, error)
	Info() Info
}

type cgroupCPU struct {
	frequency uint64
	quota     float64 // 表示 限制当前cgroup 可以使用的cpu时间 单位毫秒
	cores     uint64  // cpu数量

	preSystem uint64
	preTotal  uint64 // 每次读取 cpu时  cpu总的运行时间
	usage     uint64
}

// 这里实际 需要去读取 cgroup 文件的值 然后 /sys/fs/cgroup 中去读取对应的值
func (cpu *cgroupCPU) Usage() (uint64, error) {
	return 0, nil
}

func (cpu *cgroupCPU) Info() Info {
	return Info{
		Frequency: cpu.frequency,
		Quota:     cpu.quota,
	}
}

func Init() {
	stats = &cgroupCPU{}
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()

		for range ticker.C {
			u, err := stats.Usage()
			if err == nil && u != 0 {
				atomic.StoreUint64(&usage, u)
			}
		}
	}()
}

type State struct {
	Usage uint64 // CPU 的使用率
}

type Info struct {
	Frequency uint64
	Quota     float64
}

func ReadState(state *State) {
	state.Usage = atomic.LoadUint64(&usage)
}

func GetInfo() Info {
	return stats.Info()
}
