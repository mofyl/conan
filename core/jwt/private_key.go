package jwt

import (
	"crypto/md5"
	"sync"
	"time"
	"conan/config"
	"conan/log"
	"conan/utils"
)

const (
	CONF_KEY = "privateKey"
)

var (
	defaultConf PrivateConf
)

func init() {
	err := config.Decode(config.CONFIG_FILE_NAME, CONF_KEY, &defaultConf)

	if err != nil {
		panic("privateKey Conf init Fail" + err.Error())
	}
}

type PrivateConf struct {
	ExpTime int64 `yaml:"exp_time"`
}

type PrivateKey struct {
	key     []byte
	expTime int64
	lock    sync.RWMutex
	conf    *PrivateConf
}

func NewPrivateKey(conf *PrivateConf) (*PrivateKey, error) {

	if conf == nil {
		conf = &defaultConf
	}

	p := &PrivateKey{
		expTime: time.Now().Add(time.Duration(conf.ExpTime) * time.Second).Unix(),
		lock:    sync.RWMutex{},
	}

	key, err := p.getPrivateKey()

	if err != nil {
		return nil, err
	}
	p.key = key
	p.conf = conf
	return p, nil
}

func (p *PrivateKey) GetKey() []byte {
	// 查看Key是否过期
	p.lock.RLock()
	t := time.Now().Unix()
	if p.expTime > t {
		p.lock.RUnlock()
		return p.key
	}
	p.lock.RUnlock()

	// 需要更新key
	p.lock.Lock()
	key, err := p.getPrivateKey()

	if err != nil {
		log.Error("Get Private Key Fail Err is %s", err.Error())
		return nil
	}

	p.key = key
	p.expTime = time.Now().Add(time.Duration(p.conf.ExpTime) * time.Second).Unix()
	p.lock.Unlock()

	return p.key
}

func (p *PrivateKey) getPrivateKey() ([]byte, error) {
	uid, err := utils.GetUUID()

	if err != nil {
		return nil, err
	}
	uidB, err := uid.MarshalBinary()
	if err != nil {
		return nil, err
	}
	sumB := md5.Sum(uidB)

	return sumB[:], nil
}
