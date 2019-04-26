package client

import (
	"github.com/funlake/gopkg/jobworker"
	"sync"
	"time"
)

type SyncWorker struct {
	second time.Duration
	dispatcher *jobworker.BlockingDispatcher
}

func (sw *SyncWorker) Do (configs sync.Map){
	sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		configs.Range(func(key, value interface{}) bool {

			return true
		})
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {

	}))
}