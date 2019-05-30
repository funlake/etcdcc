package jobworker

import (
//"github.com/funlake/gopkg/utils/log"
//"github.com/funlake/gopkg/utils/log"
)

func NewSimpleJob(do func(), id func() string, wf func(dispatcher *BlockingDispatcher)) *simpleJob {
	job := &simpleJob{
		do,
		id,
		wf,
	}
	return job
}

type simpleJob struct {
	do func()
	id func() string
	wf func(dispatcher *BlockingDispatcher)
}

func (sj *simpleJob) Do() {
	sj.do()
	//log.Info("helloworld")
	//time.Sleep(2)
}
func (sj *simpleJob) Id() string {
	return sj.id()
}
func (sj *simpleJob) OnWorkerFull(dispatcher *BlockingDispatcher) {
	sj.wf(dispatcher)
	//log.Warning("Worker is full")
	//dispatcher.Put(sj)
}
