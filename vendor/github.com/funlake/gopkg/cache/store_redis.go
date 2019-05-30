package cache

import (
	"github.com/funlake/gopkg/utils/log"
	"github.com/garyburd/redigo/redis"
	"sync"
	"time"
)

func NewKvStoreRedis() *KvStoreRedis {
	return &KvStoreRedis{}
}

type KvStoreRedis struct {
	cacheSync sync.Once
	pool      *redis.Pool
}

func (sr *KvStoreRedis) Connect(dsn string, pwd string) {
	sr.cacheSync.Do(func() {
		sr.pool = &redis.Pool{
			MaxIdle:     50,
			MaxActive:   150,
			IdleTimeout: 100 * time.Second,
			Wait:        false,
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", dsn)
				if err != nil {
					log.Error(err.Error())
					return nil, err
				}
				if pwd != "" {
					if _, err := c.Do("AUTH", pwd); err != nil {
						log.Error(err.Error())
						c.Close()
						return nil, err
					}
				}
				// 选择db
				c.Do("SELECT", 0)
				log.Success("Init one redis connection")
				return c, nil
			},
		}
	})
}

func (sr *KvStoreRedis) ConnectWithTls(dsn, tls interface{}) error {
	return nil
}
func (sr *KvStoreRedis) Set(key string, val interface{}) {
	c := sr.pool.Get()
	defer c.Close()
	c.Do("SET", key, val)
}
func (sr *KvStoreRedis) Get(key string) (interface{}, error) {
	c := sr.pool.Get()
	defer c.Close()
	return redis.String(c.Do("GET", key))
}
func (sr *KvStoreRedis) HashGet(hkey string, key string) (interface{}, error) {
	if hkey == "" {
		return sr.Get(key)
	}
	c := sr.pool.Get()
	defer c.Close()
	return redis.String(c.Do("HGET", hkey, key))
}
func (sr *KvStoreRedis) HashSet(hkey string, key string, val interface{}) (interface{}, error) {
	c := sr.pool.Get()
	defer c.Close()
	return redis.String(c.Do("HSET", hkey, key, val))
}
func (sr *KvStoreRedis) GetPool() interface{} {
	return sr.pool
}
func (sr *KvStoreRedis) GetActiveCount() int {
	return sr.pool.ActiveCount()
}
