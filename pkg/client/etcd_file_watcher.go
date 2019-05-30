package client

import (
	"context"
	"etcdcc/pkg/log"
	"github.com/funlake/gopkg/jobworker"
	"os"
	"strings"
)

//EtcdFileWatcher : file storage watcher for etcd
type EtcdFileWatcher struct {
	RetrySeconds int
	StoreDir     string
	GeneralWatcher
}

//KeepEyesOnKey : Watching specific key
func (ecw *EtcdFileWatcher) KeepEyesOnKey(key string) {}

//KeepEyesOnKeyWithPrefix : Watching specific prefix
func (ecw *EtcdFileWatcher) KeepEyesOnKeyWithPrefix(module string) {
	ecw.setWorker(module)
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
