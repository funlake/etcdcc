package jobworker

import (
	//"github.com/funlake/gopkg/utils/log"
	"github.com/funlake/gopkg/utils/log"
)

type Worker struct {
	workerPool chan chan WorkerJob
	jobChannel chan WorkerJob
}

func NewWorker(workerPool chan chan WorkerJob) Worker {
	worker := Worker{
		workerPool: workerPool,
		jobChannel: make(chan WorkerJob),
	}
	return worker
}

func (w Worker) Ready() {
	go func() {
		defer func() {
			if err := recover(); err != nil {
				log.Warning("Restarting worker since panic invoke : %s", err)
				(NewWorker(w.workerPool)).Ready()
			}
		}()
		for {
			//注册进worker 池,以供dispatcher消费,传入job
			w.workerPool <- w.jobChannel
			select {
			//发现dispatcher传入了任务
			case job := <-w.jobChannel:
				job.Do()
			}
		}
	}()
}
