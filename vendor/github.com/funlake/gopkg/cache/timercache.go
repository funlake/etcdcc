package cache

type TimerCache interface {
	Flush()
	Get(hk string, k string, wheel int) (string, error)
}
