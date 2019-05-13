package client

import (
	"context"
	"etcdcc/apiserver/pkg/log"
	"github.com/funlake/gopkg/jobworker"
	"os"
	"strings"
)

type EtcdFileWatcher struct {
	RetrySeconds int
	StoreDir     string
	GeneralWatcher
}

func (ecw *EtcdFileWatcher) KeepEyesOnKey(key string) {}

func (ecw *EtcdFileWatcher) KeepEyesOnKeyWithPrefix(module string) {
	ecw.SetWorker(module)
}
func (ecw *EtcdFileWatcher) SetWorker(module string) {
	storeDir := ecw.StoreDir + "/" + module
	err := os.MkdirAll(storeDir, 0755)
	if err != nil {
		log.Error("Can not create directory for configuration files : " + err.Error())
		return
	}
	syncWorker := &SyncFileWorker{
		storeDir:     storeDir,
		shmfile:      strings.Replace(module, "/", "_", -1),
		retrySeconds: ecw.RetrySeconds,
		//Big queue shared with single syncWorker
		dispatcher: jobworker.NewBlockingDispather(1, 200),
	}
	go syncWorker.Retry()
	//Initialize all configurations under mod
	ecw.Init(module, func(k, v string) {
		syncWorker.SyncOne(k, v)
	})
	ecw.Watch(module, func(k, v string) {
		syncWorker.SyncOne(k, v)
	}, func(mk, k string, cancel context.CancelFunc) {
		//监听key == 删除key，则整个watch停止
		if mk == k {
			cancel()
		}
		syncWorker.RemoveOne(k)
	})
}
//func (ecw *EtcdFileWatcher) Init(moduleKey string, callback func(k, v string)) {
//	log.Info(fmt.Sprintf("Initialize configuration with %s", moduleKey))
//	adapter := etcd.Adapter{}
//	allKeys, err := adapter.GetMetaCacheHandler().GetStore().Get(moduleKey+"/", clientv3.WithPrefix())
//	if err == nil {
//		for _, e := range allKeys.(*clientv3.GetResponse).Kvs {
//			sk := strings.TrimPrefix(string(e.Key), moduleKey+"/")
//			//ecw.ModifyLocal(sk, string(e.Value))
//			callback(sk, string(e.Value))
//		}
//		//syncWorker.SyncAll(ecw.configs)
//	} else {
//		log.Error(err.Error())
//		return
//	}
//}
//func (ecw *EtcdFileWatcher) Watch(key string, putCallback func(k, v string), delCallBack func(mk, k string, cancel context.CancelFunc)) {
//	adapter := etcd.Adapter{}
//	ctx, cancel := context.WithCancel(context.Background())
//	log.Info(fmt.Sprintf("Watching key with %s", key))
//	//Watching mod's configurations
//	for v := range adapter.GetMetaCacheHandler().GetStore().Watch(ctx, key, clientv3.WithPrefix()) {
//		if v.Err() != nil {
//			continue
//		}
//		for _, e := range v.Events {
//			tp := fmt.Sprintf("%v", e.Type)
//			sk := strings.TrimPrefix(string(e.Kv.Key), key+"/")
//			switch tp {
//			case "PUT":
//				putCallback(sk, string(e.Kv.Value))
//			case "DELETE":
//				delCallBack(key, sk, cancel)
//			}
//		}
//	}
//	cancel()
//}

//func (ecw *EtcdFileWatcher) ModifyLocal(key, val string) {
//	ecw.configs.Store(key, val)
//}
