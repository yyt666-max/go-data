package cache_redis

import (
	"context"
	"log"
	"time"

	"github.com/eolinker/go-common/autowire"
	"github.com/eolinker/go-common/cache"
	"github.com/eolinker/go-common/cftool"
)

type RedisConfig struct {
	UserName   string   `yaml:"user_name"`
	Password   string   `yaml:"password"`
	Addr       []string `yaml:"addr"`
	Prefix     string   `yaml:"prefix"`
	Cluster    string   `yaml:"cluster"`
	MasterName string   `yaml:"master_name"`
	DB         int      `yaml:"db"`
}

type redisInit struct {
	conf *RedisConfig `autowired:""`
}

func (r *redisInit) OnComplete() {

	client := SimpleCluster(Option{
		Addrs:      r.conf.Addr,
		MasterName: r.conf.MasterName,
		Username:   r.conf.UserName,
		Password:   r.conf.Password,
		DB:         r.conf.DB,
	})

	timeout, cancelFunc := context.WithTimeout(context.Background(), time.Second*10)
	defer cancelFunc()
	if err := client.Ping(timeout).Err(); err != nil {
		_ = client.Close()

		log.Fatalf("ping redis %v error:%s", r.conf.Addr, err.Error())
	}

	iCommonCache := newCommonCache(client, r.conf.Prefix)

	autowire.Inject[cache.ICommonCache](iCommonCache)

}

func init() {
	cftool.Register[RedisConfig]("redis")
	autowire.Autowired(new(redisInit))
}
