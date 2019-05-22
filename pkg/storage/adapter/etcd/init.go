package etcd

import (
	"etcdcc/apiserver/pkg/log"
	"github.com/coreos/etcd/pkg/transport"
	_cache "github.com/funlake/gopkg/cache"
	"sync"
)

var (
	EtcdCache   *_cache.TimerCacheEtcd
	adapterOnce sync.Once
)

type Adapter struct{}

func (e Adapter) Connect(hosts, c, k, ca, sn string) {
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
		EtcdCache = _cache.NewTimerCacheEtcd()
		etcdStore := _cache.NewKvStoreEtcd()
		err = etcdStore.ConnectWithTls(hosts, tlsConfig)
		if err != nil {
			log.Fatal("Etcd connected failure : " + err.Error())
		}
		EtcdCache.SetStore(etcdStore)
	})
}
func (e Adapter) GetMetaCacheHandler() *_cache.TimerCacheEtcd {
	return EtcdCache
}
