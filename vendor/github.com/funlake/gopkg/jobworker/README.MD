##### Intro
任务调度器
###### Install
go get github.com/funlake/gopkg/jobworker
###### Import
import github.com/funlake/gopkg/jobworker
###### Use
```
//10个worker处理size = 100的job queue
dispatcher := jobworker.NewDispather(10,100)
//simpleJob实现WorkerJob接口
/*
type WorkerJob interface{
  	Do() //运行任务，内部如果加了协程则为非阻塞型任务调度
  	Id() string //名称
  	OnWorkerFull(dispatcher *Dispatcher) // worker已满事件
  }
*/
dispatcher.Put(simpleJob{})
//http请求构建可查看dispatche_test.go文件
```
