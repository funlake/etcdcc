package etcd

import (
	"github.com/coreos/etcd/pkg/transport"
	"github.com/funlake/etcdcc/pkg/log"
	_cache "github.com/funlake/gopkg/cache"
	"sync"
)

var (
	etcdCache   *_cache.TimerCacheEtcd
	adapterOnce sync.Once
)

//Adapter : Etcd adapter for dao layer
//type Adapter struct{}

//Connect to etcd server
func Connect(hosts, c, k, ca, sn string) {
	adapterOnce.Do(func() {
		tlsInfo := transport.TLSInfo{
			CertFile:      c,
			KeyFile:       k,
			TrustedCAFile: ca,
			ServerName:    sn,
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			log.Fatal(err.Error())
		}
		etcdCache = _cache.NewTimerCacheEtcd()
		etcdStore := _cache.NewKvStoreEtcd()
		err = etcdStore.ConnectWithTls(hosts, tlsConfig)
		if err != nil {
			log.Fatal("Etcd connected failure : " + err.Error())
		}
		etcdCache.SetStore(etcdStore)
	})
}

//GetMetaCacheHandler : Export the cache instance
func GetMetaCacheHandler() *_cache.TimerCacheEtcd {
	return etcdCache
}
