package cache

//本地介于redis之上加一层缓存
//自动定时更新缓存
//todo : lru 实现
import (
	"github.com/funlake/gopkg/timer"
	"github.com/funlake/gopkg/utils/log"
	"strings"
	"sync"
)

type TimerCacheRedis struct {
	//mu         sync.Mutex
	store *KvStoreRedis
	//local      map[string]string
	ticker     *timer.Ticker
	emptyCount map[string]int
	mcache     sync.Map
}

func NewTimerCacheRedis() *TimerCacheRedis {
	return &TimerCacheRedis{ /*local: make(map[string]string), */ ticker: timer.NewTicker(), emptyCount: make(map[string]int)}
}
func (tc *TimerCacheRedis) Flush() {
	//tc.mu.Lock()
	//defer tc.mu.Unlock()
	//for k := range tc.local {
	//	delete(tc.local, k)
	//	//ticker.Stop(k)
	//}
	tc.mcache.Range(func(key, value interface{}) bool {
		tc.mcache.Delete(key)
		return true
	})
}
func (tc *TimerCacheRedis) SetStore(store *KvStoreRedis) {
	tc.store = store
}
func (tc *TimerCacheRedis) GetStore() *KvStoreRedis {
	return tc.store
}
func (tc *TimerCacheRedis) Get(hk string, k string, wheel int) (string, error) {
	//tc.mu.Lock()
	//defer tc.mu.Unlock()
	localCacheKey := hk + "_" + k
	log.Info(localCacheKey)
	machsVal, has := tc.mcache.Load(localCacheKey)
	//if _, ok := tc.local[localCacheKey]; ok {
	if has {
		return machsVal.(string), nil
		//return tc.local[localCacheKey], nil
	} else {
		//log.Info("Access redis for setting : %s_%s",hk,k)
		v, err := tc.store.HashGet(hk, k)
		if err == nil {
			tc.ticker.Set(wheel, localCacheKey, func() {
				//log.Info("每%d秒定时检查%s",wheel,localCacheKey)
				v, err := tc.store.HashGet(hk, k)
				//假如redis服务器挂了,得保留之前的本地缓存值
				if err != nil {
					if strings.Contains(err.Error(), "nil returned") {
						if _, ok := tc.emptyCount[localCacheKey]; ok {
							tc.emptyCount[localCacheKey]++
						} else {
							tc.emptyCount[localCacheKey] = 1
						}
						//如果值返回为空，且超过3次，则停掉定时更新器(测试环境发现偶尔会有异常情况下的空值返回)
						if tc.emptyCount[localCacheKey] > 3 {
							log.Error("Empty value deteced(%d s) at least 3 times : %s,remove ticker run: error:%s", wheel, localCacheKey, err.Error())
							tc.ticker.Stop(wheel, localCacheKey)
							delete(tc.emptyCount, localCacheKey)
						}
						//赋空值,如要情况缓存，可调用/api-cleancache接口
						//tc.local[localCacheKey] = v.(string)
						tc.mcache.Store(localCacheKey, v.(string))
						//delete(tc.local,localCacheKey)
					} else {
						//发生redis连接故障，则继续保持旧有缓存
						log.Error("Redis seems has gone,we do not clear cache if redis is down to keep gateway service's running")
					}
					//delete(tc.local, localCacheKey)
				} else {
					if _, ok := tc.emptyCount[localCacheKey]; ok {
						//成功清楚空值计数
						delete(tc.emptyCount, localCacheKey)
					}
					//log.Warning("更新数据%s : %s",localCacheKey,v.(string))
					//tc.local[localCacheKey] = v.(string)
					tc.mcache.Store(localCacheKey, v.(string))
				}
			})
			return v.(string), nil
			//return tc.local[localCacheKey], nil
		} else {
			log.Warning("%s : %s", localCacheKey, err.Error())
			//防止redis被刷
			//如要情况缓存，可调用/api-cleancache接口
			//tc.local[localCacheKey] = ""
			tc.mcache.Store(localCacheKey, "")
			//return "",err
		}
	}
	return machsVal.(string), nil
	//log.Warning("waht the fuck?")
	//return tc.local[localCacheKey], nil
}
