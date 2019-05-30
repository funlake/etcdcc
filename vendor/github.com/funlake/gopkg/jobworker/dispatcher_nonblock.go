package jobworker

import (
	"github.com/funlake/gopkg/utils"
	"github.com/funlake/gopkg/utils/log"
	"time"
)

type NonBlockingDispatcher struct {
	workerPool chan chan WorkerJob
	jobQueue   chan WorkerJob
	stop       chan struct{}
	isStop     bool
}

func NewNonBlockingDispather(maxWorker int, queueSize int) *NonBlockingDispatcher {
	dispatcher := &NonBlockingDispatcher{}
	//流水线长度
	dispatcher.jobQueue = make(chan WorkerJob, queueSize)
	//流水线工人数量
	dispatcher.workerPool = make(chan chan WorkerJob, maxWorker)
	dispatcher.stop = make(chan struct{})
	dispatcher.isStop = false
	dispatcher.Run(maxWorker)
	//稍微等下worker启动
	time.Sleep(time.Nanosecond * 10)
	return dispatcher
}
func (d *NonBlockingDispatcher) Put(job WorkerJob) bool {
	if d.isStop {
		return false
	}
	select {
	case <-d.stop:
		close(d.jobQueue)
		//close(d.workerPool)
		d.isStop = true
		return false
	case d.jobQueue <- job:
		return true
	default:
		return false
		//log.Error("job队列已满")
	}
	return false
}
func (d *NonBlockingDispatcher) Run(maxWorker int) {
	for i := 0; i < maxWorker; i++ {
		worker := NewWorker(d.workerPool)
		worker.Ready()
	}
	utils.WrapGo(func() {
		d.Ready()
	}, "dispather run")
}

//需做worker繁忙超时处理机制,否则会因为worker过于繁忙，
//症状是任务快速加入，然而worker繁忙导致管道堵塞，并发时，如果worker数量远小于并发，
//如不做超时处理,则后续请求会在<-d.workerPool处等待很长时间，而真正进入转发阶段的消耗时间却并不多
//所以后端服务无压力，而请求却堵在网关处!
func (d *NonBlockingDispatcher) Ready() {
stop:
	for {
		select {
		case job, open := <-d.jobQueue:
			if !open {
				log.Error("Closing workerPool")
				close(d.workerPool)
				break stop
			}
			select {
			case jobChan := <-d.workerPool:
				jobChan <- job
			default:
				job.(WorkerNonBlockingJob).OnWorkerFull(d)
				//job.OnWorkerFull(d)
			}
		}
	}
	log.Info("Worker pool has been shutted down")
}
func (d *NonBlockingDispatcher) Stop() {
	d.stop <- struct{}{}
}
func (d *NonBlockingDispatcher) GetActiveWorkers() int {
	return len(d.workerPool)
}

func (d *NonBlockingDispatcher) GetActiveQueue() int {
	return len(d.jobQueue)
}
func (d *NonBlockingDispatcher) StopStatus() bool {
	return d.isStop
}
