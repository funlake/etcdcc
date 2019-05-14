package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"etcdcc/apiserver/pkg/log"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"github.com/tidwall/gjson"
	"net"
	"os"
	"strings"
	"sync"
)

const (
	TypeYaml = "yaml"
	TypeJson = "json"
	TypeToml = "toml"
	TypeProp = "prop"
)

type EtcdUdsWatcher struct {
	GeneralWatcher
	rawConfig sync.Map
}

func (euw *EtcdUdsWatcher) KeepEyesOnKey(key string) {}

//1. watch etcd
//2. serve uds
func (euw *EtcdUdsWatcher) KeepEyesOnKeyWithPrefix(prefix string) {
	euw.Init(prefix, func(k, v string) {
		euw.SaveLocal(k, v)
	})
	euw.Watch(prefix, func(k, v string) {
		euw.SaveLocal(k, v)
	}, func(mk, k string, cancel context.CancelFunc) {
		if mk == k {
			cancel()
		}
		euw.rawConfig.Delete(k)
	})
}
func (euw *EtcdUdsWatcher) SaveLocal(k, v string) {
	r, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		log.Error("Base64 decode error:" + err.Error())
	} else {
		//transform yaml/toml to json for easy handling
		if strings.HasPrefix(k, TypeYaml+"/") {
			r, err = yaml.YAMLToJSON(r)
			if err != nil {
				log.Error("Yaml to json error :" + err.Error())
			}
		}
		if strings.HasPrefix(k, TypeToml+"/") {
			var tm interface{}
			_, err = toml.Decode(string(r), &tm)
			if err != nil {
				log.Error("Toml decode error :" + err.Error())
			}
			r, err = json.Marshal(tm)
		}
		euw.rawConfig.Store(k, string(r))
	}
}
func (euw *EtcdUdsWatcher) ServeSocket(sockFile string) {
	_ = os.Remove(sockFile)
	ln, err := net.Listen("unix", sockFile)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer ln.Close()
	log.Info("Unix domain socket listening on " + sockFile)
	for {
		fd, err := ln.Accept()
		if err != nil {
			log.Error("Accept error: " + err.Error())
		} else {
			log.Info("New client coming")
		}
		go euw.handle(fd)
	}
}

func (euw *EtcdUdsWatcher) handle(fd net.Conn) {
	for {
		buf := make([]byte, 512)
		n, err := fd.Read(buf)
		if err != nil {
			return
		}
		msg := string(buf[0:n])
		cmd := strings.SplitN(msg, " ", 3)
		switch strings.ToLower(cmd[0]) {
		case "get":
			r, ok := euw.rawConfig.Load(cmd[1])
			if ok {
				val := ""
				if len(cmd) > 2 {
					ts := strings.SplitN(cmd[1], "/", 2)
					val = euw.getSpecifyKey(r.(string), cmd[2], ts[0])
					if val != "" {
						_, err = fd.Write([]byte(val))
					} else {
						err = errors.New("Not found")
					}
				} else {
					_, err = fd.Write([]byte(r.(string)))
				}
				if err != nil {
					_, _ = fd.Write([]byte(err.Error()))
					log.Error("Uds response error:" + err.Error())
				}
			} else {
				_, _ = fd.Write([]byte("Not found " + cmd[1]))
			}
		case "raw":
			r, ok := euw.rawConfig.Load(cmd[1])
			var c []byte
			if ok {
				ts := strings.SplitN(cmd[1], "/", 2)
				switch ts[0] {
				case TypeYaml:
					c, _ = yaml.JSONToYAML([]byte(r.(string)))
				case TypeToml:
					var b bytes.Buffer
					var m interface{}
					_ = json.Unmarshal([]byte(r.(string)), &m)
					_ = toml.NewEncoder(&b).Encode(m)
					c = b.Bytes()
				default:
					c = []byte(r.(string))
				}
				_, err = fd.Write(c)
			} else {
				_, _ = fd.Write([]byte("No specify configuration for " + cmd[1]))
			}
		default:
			log.Error("Unknown command:[" + msg + "]")
		}
	}
}

func (euw *EtcdUdsWatcher) getSpecifyKey(raw, k, t string) string {
	switch t {
	case TypeJson, TypeToml, TypeYaml:
		return gjson.Get(raw, k).String()
	}
	return raw
}
