package client

import "time"

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
