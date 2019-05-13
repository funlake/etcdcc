package client

import (
	"context"
	"encoding/base64"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type EtcdUdsWatcher struct {
	GeneralWatcher
	localConfig sync.Map
}
func(euw *EtcdUdsWatcher) KeepEyesOnKey (key string){}
//1. watch etcd
//2. serve uds
func(euw *EtcdUdsWatcher) KeepEyesOnKeyWithPrefix(module string){
	euw.Init(module, func(k, v string) {
		euw.localConfig.Store(k,v)
	})
	euw.Watch(module, func(k, v string) {
		euw.localConfig.Store(k,v)
	}, func(mk, k string, cancel context.CancelFunc) {
		if mk == k {
			cancel()
		}
		euw.localConfig.Delete(k)
	})
}
func (euw *EtcdUdsWatcher) Serve()  {
	_ = os.Remove("/dev/shm/etcdcc.sock")
	ln, err := net.Listen("unix", "/dev/shm/etcdcc.sock")

	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()
	log.Info("Unix domain socket listening on /dev/shm/etcdcc.sock")
	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Error("Accept error: "+err.Error())
		} else {
			log.Info("New client coming")
		}
		go euw.handle(fd)
	}
}

func (euw *EtcdUdsWatcher) handle(fd net.Conn)  {
	for {
		buf := make([]byte, 512)
		n, err := fd.Read(buf)
		if err != nil {
			return
		}
		msg := string(buf[0:n])
		cmd := strings.SplitN(msg, " ", 2)
		switch strings.ToLower(cmd[0]) {
		case "get":
			v,ok := euw.localConfig.Load(cmd[1])
			if ok {
				r,_ :=  base64.StdEncoding.DecodeString(v.(string))
				_, err = fd.Write([]byte(fmt.Sprintf("Configuration %s",r)))
				if err != nil {
					log.Error("Uds response error:" + err.Error())
				}
			} else {
				_, _ = fd.Write([]byte("No specify configuration with key:" + cmd[1]))
			}
		default:
			log.Error("Unknown command:[" + msg + "]")
		}
	}
}
