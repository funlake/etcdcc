package dao

import (
	"github.com/funlake/etcdcc/pkg/storage/adapter/etcd"
	"testing"
)

var metaTestModel = MetaCache{}

func initConnect() {
	etcd.Connect("https://127.0.0.1:2479", "/keys/client.pem", "/keys/client-key.pem", "/keys/ca.pem", "")
}

//func TestEtcdService_Delete(t *testing.T) {
//	initConnect()
//	r, err := metaTestModel.Delete("dev/act/conf/nginx")
//	if err != nil {
//		t.Log(err.Error())
//	} else {
//		t.Log(r)
//	}
//}

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
