package client

// @Summary
// 1. 同步器用job worker异步阻塞模式，只开一个工作协程，避免写锁,协程panic重启由类库自动实现
// 2. 定时重试
// 3. 只兼容linux,一些命令用linux原生执行(通用后期再考虑)
type Watcher interface {
	KeepEyesOnKey(key string)
	KeepEyesOnKeyWithPrefix(key string)
	ModifyLocal(key, val string)
}
