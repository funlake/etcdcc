package client

import (
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/funlake/gopkg/jobworker"
	"sync"
)

type SyncWorker struct {
	second int
	dispatcher *jobworker.BlockingDispatcher
}

func (sw *SyncWorker) Do (configs sync.Map){
	sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		configs.Range(func(key, value interface{}) bool {
			log.Info(fmt.Sprintf("%s:%s",key,value))
			return true
		})
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {
		//never invoke..
	}))
}