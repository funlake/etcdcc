package client

import "time"

//Interface of worker
type Worker interface {
	RemoveOne(key interface{})
	SyncOne(key, value interface{})
	Retry()
}

type failConfig struct {
	t time.Time
	k interface{}
	v interface{}
}
