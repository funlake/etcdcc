package jobworker

type WorkerJob interface {
	Do()
	Id() string
	//OnWorkerFull(dispatcher Dispatcher)
}
type WorkerNonBlockingJob interface {
	OnWorkerFull(dispatcher *NonBlockingDispatcher)
}
type WorkerBlockingJob interface {
	OnWorkerFull(dispatcher *BlockingDispatcher)
}
