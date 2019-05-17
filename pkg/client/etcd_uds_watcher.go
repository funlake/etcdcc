package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"etcdcc/apiserver/pkg/log"
	"github.com/BurntSushi/toml"
	"github.com/ghodss/yaml"
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

var readBuffer = make([]byte, 1024)

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
		r, err = euw.jsonEncode(r, k)
		euw.rawConfig.Store(k, string(r))
	}
}
func (euw *EtcdUdsWatcher) jsonEncode(r []byte, prefix string) ([]byte, error) {
	var err error
	if strings.HasPrefix(prefix, TypeYaml+"/") {
		r, err = yaml.YAMLToJSON(r)
		if err != nil {
			log.Error("Yaml to json error :" + err.Error())
		}
	}
	if strings.HasPrefix(prefix, TypeToml+"/") {
		var tm interface{}
		_, err = toml.Decode(string(r), &tm)
		if err != nil {
			log.Error("Toml decode error :" + err.Error())
		}
		r, err = json.Marshal(tm)
	}
	return r, err
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
		cmd, err := euw.getCmd(fd)
		if err != nil {
			log.Error("Get cmd error:" + err.Error())
			_, _ = fd.Write([]byte("fail," + err.Error()))
			return
		}
		switch strings.ToLower(cmd[0]) {
		case "get":
			val, err := euw.find(cmd)
			if err != nil {
				_, err = fd.Write([]byte("fail," + err.Error()))
			} else {
				_, err = fd.Write([]byte("ok," + val))
			}
			if err != nil {
				_, _ = fd.Write([]byte("fail," + err.Error()))
				log.Error("Uds response error:" + err.Error())
			}
		//sometimes need to check raw format things
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
			log.Error("Unknown command:[" + strings.Join(cmd, " ") + "]")
		}
	}
}
func (euw *EtcdUdsWatcher) getCmd(fd net.Conn) ([]string, error) {
	n, err := fd.Read(readBuffer)
	if err != nil {
		return []string{}, err
	}
	msg := string(readBuffer[0:n])
	cmd := strings.SplitN(msg, " ", 3)
	return cmd, nil
}
func (euw *EtcdUdsWatcher) find(cmd []string) (string, error) {
	r, ok := euw.rawConfig.Load(cmd[1])
	if ok {
		if len(cmd) > 2 {
			val := euw.getSpecifyKey(r.(string), cmd)
			if val != "" {
				return val, nil
			}
		} else {
			return r.(string), nil
		}
	}
	return "", errors.New("NotFound")
}
func (euw *EtcdUdsWatcher) getSpecifyKey(raw string, cmd []string) string {
	t := strings.SplitN(cmd[1], "/", 2)
	switch t[0] {
	case TypeJson, TypeToml, TypeYaml:
		return gjson.Get(raw, cmd[2]).String()
	}
	return raw
}
