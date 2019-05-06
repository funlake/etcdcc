package models

import (
	"context"
	"etcdcc/apiserver/pkg/dao/adapter/etcd"
	"fmt"
	"github.com/coreos/etcd/clientv3"
	"testing"
)

var metaTestAdapter = etcd.Adapter{}
var metaTestModel = MetaCache{}

func initConnect() {
	metaTestAdapter.Connect("https://127.0.0.1:2479", "/keys/ca.pem", "/keys/ca-key.pem", "/keys/ca.crt", "etcchebao")
}
func TestEtcdService_Get(t *testing.T) {
	initConnect()
	r, err := metaTestModel.Get("dev/act/act.conf")
	if err != nil {
		t.Log(err.Error())
	} else {
		t.Log(r)
		//Watch()
	}
}

func Watch() {
	ctx, cancel := context.WithCancel(context.Background())
	for v := range metaTestAdapter.GetMetaCacheHandler().GetStore().Watch(ctx, "dev/gateway", clientv3.WithPrefix()) {
		if v.Err() != nil {
			cancel()
		}
		for _, e := range v.Events {
			tp := fmt.Sprintf("%v", e.Type)
			fmt.Println(tp, string(e.Kv.Key), string(e.Kv.Value))
		}
	}
}

//func TestEtcdService_Put(t *testing.T) {
//	initConnect()
//	_, err := metaTestModel.Put("foo/hello", "hello world")
//	if err != nil {
//		t.Error(err.Error())
//	} else {
//		t.Log("Set ok!")
//	}
//}
func BenchmarkEtcdService_Get(b *testing.B) {
	initConnect()
	b.SetParallelism(10)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = metaTestModel.Get("")
		}
	})
}
func BenchmarkEtcdService_Put(b *testing.B) {
	initConnect()
	b.SetParallelism(100)
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, _ = metaTestModel.Put("foo", "www")
		}
	})
}
