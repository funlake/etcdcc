package cache

type KvStore interface {
	Connect(dsn, pwd string)
	ConnectWithTls(dsn, tls interface{})
	Get(key string) (interface{}, error)
	Set(key string, val interface{})
	GetPool() interface{}
}
