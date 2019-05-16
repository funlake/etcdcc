package services

import (
	"etcdcc/apiserver/pkg/storage/adapter/etcd"
	"etcdcc/apiserver/pkg/storage/adapter/mysql"
	"etcdcc/apiserver/pkg/dto"
	"testing"
)
func initConfigTest(){
	mysql.Adapter.Connect(mysql.Adapter{})
	etcd.Adapter.Connect(etcd.Adapter{},"https://127.0.0.1:2479","/keys/ca.pem","/keys/ca-key.pem","/keys/ca.crt","etcchebao")
}
func TestConfig_Crud(t *testing.T) {
	initConfigTest()
	fake := &dto.ConfigAddDto{
		Env: "dev",
		Mod: "act",
		Key: "goodboy",
		Val: "www-dev2",
		Type: "json",
	}
	cs := Config{}
	id ,err := cs.Create(fake)
	if err != nil{
		t.Fatal(err.Error())
	}
	fake2 := &dto.ConfigEditDto{
		Id: int(id),
		Env: "dev",
		Mod: "act",
		Key: "goodboy",
		Val: `{"a":"b",
		"c":4}`,
		Type: "json",
	}
	err = cs.Update(fake2)
	if err != nil{
		t.Fatal(err.Error())
	}
	err = cs.Delete(&dto.ConfigDelDto{
		Id: int(id),
	})
	if err != nil{
		t.Fatal(err.Error())
	}
}

