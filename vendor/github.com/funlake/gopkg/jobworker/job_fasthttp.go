package jobworker

import (
	"errors"
	"github.com/funlake/gopkg/utils/log"
	"github.com/valyala/fasthttp"
	"sync"
	"time"
)

//var fastHttpResChan chan chan FastHttpProxyJobResponse
var fasthttpOnce = &sync.Once{}

type FastHttpProxyJobResponse struct {
	Response *fasthttp.Response
	Error    error
	Dur      time.Duration
}

var fastHttpProxyJobPool = sync.Pool{
	New: func() interface{} {
		return &fastHttpProxyJob{}
	},
}

//func initFastHttpProxyResChan(chanSize int){
//	fasthttpOnce.Do(func() {
//		log.Success("Res fasthttp channel size:%d",chanSize)
//		fastHttpResChan = make(chan chan FastHttpProxyJobResponse,chanSize)
//		for i:=0;i<chanSize;i++{
//			fastHttpResChan <- make(chan FastHttpProxyJobResponse)
//		}
//	})
//}
func NewFastHttpProxyJob(transport *fasthttp.Client, q *fasthttp.Request, id string) *fastHttpProxyJob {
	//initFastHttpProxyResChan(rcsize)
	job := fastHttpProxyJobPool.Get().(*fastHttpProxyJob)
	job.q = q
	job.m = id
	job.t = transport
	job.r = make(chan FastHttpProxyJobResponse)
	//job := &fastHttpProxyJob{q:q,m:id,t:transport}
	//job.setResChan()
	return job
}

type fastHttpProxyJob struct {
	q *fasthttp.Request
	r chan FastHttpProxyJobResponse
	m string
	t *fasthttp.Client
}

func (job *fastHttpProxyJob) Id() string {
	return job.m
}
func (job *fastHttpProxyJob) OnWorkerFull(dispatcher *NonBlockingDispatcher) {
	log.Error("worker 繁忙")
	job.r <- FastHttpProxyJobResponse{nil, errors.New("worker繁忙"), 0}
}
func (job *fastHttpProxyJob) Do() {
	//now := time.Now()
	resp := fasthttp.AcquireResponse()
	err := job.t.Do(job.q, resp)
	//1.加协程可快速释放worker,worker不论转发是否成功就回列
	// 优点:处理速度快,保证worker快速回,保证可用worker数
	// 缺点:假如后端高并发处理能力不足,则会造成雪崩效应，后端已经处理不过来了，这边还是持续由worker发送请求
	//2.不加协程可由worker数量+请求处理速度,利用workder堵塞,可适当保护后端服务,提升流量爆发期间的tps(后续请求会马上返回worker已满),当然worker数量的配置得和处理能力相当
	// 优点:量化处理能力，如果后端处理得过来，则增加worker数量，否则减少
	// 缺点:高并发会因为超时处理过多,job.R消费等待过长会发生worker不够用的情况
	//目前选择:2,处理能力可通过调整worker size大小来决定
	//TODO : 更多测试
	//go func(res *http.Response,err error) {
	job.r <- FastHttpProxyJobResponse{resp, err, 0}
	//time.Sleep(time.Second * 3)
	//log.Info("%s -> fasthttp请求时间消耗 : %s",job.Id(),time.Since(now))
	//}()
}
func (job *fastHttpProxyJob) GetResChan() chan FastHttpProxyJobResponse {
	return job.r
}
func (job *fastHttpProxyJob) Release() {
	<-job.r
	//job.putResChan()
	//httpProxyJobPool.Put(job)
}

//func (job *fastHttpProxyJob) setResChan(){
//	job.r = <-fastHttpResChan
//}
//
//func (job *fastHttpProxyJob) putResChan(){
//	fastHttpResChan <- job.r
//}
