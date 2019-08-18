package client

import (
	"context"
	"github.com/funlake/etcdcc/pkg/log"
	"github.com/funlake/gopkg/cache"
	"github.com/funlake/gopkg/jobworker"
	"os"
	"strings"
)

//EtcdFileWatcher : file storage watcher for etcd
type EtcdFileWatcher struct {
	RetrySeconds int
	StoreDir     string
	GeneralWatcher
	Tc *cache.TimerCacheEtcd
}

//KeepEyesOnKey : Watching specific key
func (ecw *EtcdFileWatcher) KeepEyesOnKey(key string) {}

//KeepEyesOnKeyWithPrefix : Watching specific prefix
func (ecw *EtcdFileWatcher) KeepEyesOnKeyWithPrefix(module string) {
	ecw.setWorker(module)
}
func (ecw *EtcdFileWatcher) SaveLocal(k, v string) {
	//
}
func (ecw *EtcdFileWatcher) setWorker(module string) {
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
	ecw.Init(ecw.Tc, module, func(k, v string) {
		syncWorker.SyncOne(k, v)
	})
	ecw.Watch(ecw.Tc, module, func(k, v string) {
		syncWorker.SyncOne(k, v)
	}, func(mk, k string, cancel context.CancelFunc) {
		//监听key == 删除key，则整个watch停止
		if mk == k {
			cancel()
		}
		syncWorker.RemoveOne(k)
	})
}

func (ecw *EtcdFileWatcher) Find(cmd []string) (string, error) {
	return "", nil
}
