package cache

import (
	"context"
	"errors"
	"fmt"
	cv3 "github.com/coreos/etcd/clientv3"
	"github.com/funlake/gopkg/utils"
	"sync"
)

func NewTimerCacheEtcd() *TimerCacheEtcd {
	return &TimerCacheEtcd{}
}

type TimerCacheEtcd struct {
	store *KvStoreEtcd
	local sync.Map
}

func (tc *TimerCacheEtcd) GetStore() *KvStoreEtcd {
	return tc.store
}
func (tc *TimerCacheEtcd) Get(hk string, k string, wheel int) (string, error) {
	var rv string
	rk := k
	if hk != "" {
		rk = hk + "/" + k
	}

	val, ok := tc.local.Load(rk)
	if !ok {
		resp, err := tc.store.Get(rk)
		if err == nil {
			for _, e := range resp.(*cv3.GetResponse).Kvs {
				rv = string(e.Value)
			}
		} else {
			tc.local.Store(rk, "")
			return "", err
		}
		if rv != "" {
			utils.WrapGo(func() {
				tc.Watch(rk)
			}, fmt.Sprintf("watch-key-%s", rk))
			tc.local.Store(rk, rv)
			return rv, nil
		} else {
			tc.local.Store(rk, "")
			return "", errors.New("Value Not set")
		}
	}
	rv = val.(string)
	return rv, nil
}
func (tc *TimerCacheEtcd) Flush() {
	tc.local.Range(func(key, value interface{}) bool {
		tc.local.Delete(key)
		return true
	})
}
func (tc *TimerCacheEtcd) SetStore(store *KvStoreEtcd) {
	tc.store = store
}
func (tc *TimerCacheEtcd) Watch(key string) {
	ctx, cancel := context.WithCancel(context.Background())
	wc := tc.store.Watch(ctx, key)
	for v := range wc {
		if v.Err() != nil {
			panic(v.Err().Error())
		}
		for _, e := range v.Events {
			tp := fmt.Sprintf("%v", e.Type)
			switch tp {
			case "DELETE":
				tc.local.Delete(key)
				cancel()
				break
			case "PUT":
				tc.local.Store(key, string(e.Kv.Value))
				break
			}
		}
	}
}
