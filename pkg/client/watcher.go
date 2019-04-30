package client

// todo
// 1. 开一个定时器,定时同步最新配置,做成启动配置
// 2. 同步器用job worker异步阻塞模式，只开一个工作协程，避免写锁
// 3. ...
type Watcher interface {
	KeepEyesOnKey(key string)
	KeepEyesOnKeyWithPrefix(key string, prefix interface{})
	ModifyLocal(key, val string)
}
