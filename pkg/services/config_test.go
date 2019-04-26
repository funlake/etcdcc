package services

import (
	"etcdcc/apiserver/pkg/dao/etcd"
	"etcdcc/apiserver/pkg/dao/mysql"
	"etcdcc/apiserver/pkg/dto"
	"testing"
)
func initConfigTest(){
	mysql.Adapter.Connect(mysql.Adapter{})
	etcd.Adapter.Connect(etcd.Adapter{},"https://127.0.0.1:2479","/keys/ca.pem","/keys/ca-key.pem","/keys/ca.crt","etcchebao")
}
func TestConfig_Create(t *testing.T) {
	initConfigTest()
	fake := &dto.ConfigAddDto{
		Env: "dev",
		Mod: "gateway",
		Key: "proxy/access_token",
		Val: "www-dev2",
	}
	cs := Config{}
	_ ,err := cs.Create(fake)
	if err != nil{
		t.Error(err.Error())
	}
}
func TestConfig_Update(t *testing.T) {
	fake := &dto.ConfigEditDto{
		Id: 10,
		Env: "dev",
		Mod: "gateway",
		Key: "proxy/access_token",
		Val: "www-dev14",
	}
	cs := Config{}
	err := cs.Update(fake)
	if err != nil{
		t.Error(err.Error())
	}
}
func TestConfig_Delete(t *testing.T) {
	//initConfigTest()
	fake := &dto.ConfigDelDto{
		Id: 5,
	}
	cs := Config{}
	err := cs.Delete(fake)
	if err != nil{
		t.Log(err.Error())
	}
}

