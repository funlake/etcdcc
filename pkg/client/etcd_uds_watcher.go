package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"etcdcc/pkg/log"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/funlake/gopkg/cache"
	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
	"github.com/zieckey/goini"
	"net"
	"os"
	"strings"
	"sync"
)





//EtcdUdsWatcher : Unix domain socket watcher for etcd
type EtcdUdsWatcher struct {
	GeneralWatcher
	rawConfig sync.Map
	Tc *cache.TimerCacheEtcd
}

//KeepEyesOnKey : Watching specific key
func (euw *EtcdUdsWatcher) KeepEyesOnKey(key string) {}

//KeepEyesOnKeyWithPrefix : Watch etcd with prefix
func (euw *EtcdUdsWatcher) KeepEyesOnKeyWithPrefix(prefix string) {
	euw.Init(euw.Tc,prefix, func(k, v string) {
		euw.saveLocal(k, v)
	})
	euw.Watch(euw.Tc,prefix, func(k, v string) {
		euw.saveLocal(k, v)
	}, func(mk, k string, cancel context.CancelFunc) {
		if mk == k {
			cancel()
		}
		euw.rawConfig.Delete(k)
	})
}
func (euw *EtcdUdsWatcher) saveLocal(k, v string) {
	r, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		log.Error(fmt.Sprintf("Base64 decode error with key %s : %s:" ,k,err.Error()))
	} else {
		//transform yaml/toml to json for easy handling
		r, err = euw.jsonEncode(r, k)
		log.Debug(fmt.Sprintf("saving data %s -> %s",k,string(r)))
		euw.rawConfig.Store(k, string(r))
	}
}
func (euw *EtcdUdsWatcher) jsonEncode(r []byte, prefix string) ([]byte, error) {
	var (
		err error
	)
	if strings.HasPrefix(prefix, typeYaml+"/") {
		r, err = yaml.YAMLToJSON(r)
		if err != nil {
			log.Error("Yaml to json error :" + err.Error())
		}
	}
	if strings.HasPrefix(prefix, typeToml+"/") {
		var tm interface{}
		_, err = toml.Decode(string(r), &tm)
		if err != nil {
			log.Error("Toml decode error :" + err.Error())
		}
		r, err = json.Marshal(tm)
	}
	if strings.HasPrefix(prefix, typeProp+"/") {
		ini := goini.New()
		err = ini.Parse(r, "\n", "=")
		if err == nil {
			r, err = json.Marshal(ini.GetAll())
		}
	}
	return r, err
}

//ServeSocket : Serve unix socket for applications
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
			return
		}
		switch strings.ToLower(cmd[0]) {
		case "get":
			val, err := euw.Find(cmd)
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
				case typeYaml:
					c, _ = yaml.JSONToYAML([]byte(r.(string)))
				case typeToml:
					var b bytes.Buffer
					var m interface{}
					_ = json.Unmarshal([]byte(r.(string)), &m)
					_ = toml.NewEncoder(&b).Encode(m)
					c = b.Bytes()
				default:
					c = []byte(r.(string))
				}
				_, _ = fd.Write(c)
			} else {
				_, _ = fd.Write([]byte("No specify configuration for " + cmd[1]))
			}
		default:
			log.Error("Unknown command:[" + strings.Join(cmd, " ") + "]")
		}
	}
}
func (euw *EtcdUdsWatcher) getCmd(fd net.Conn) ([]string, error) {
	readBuffer := make([]byte, 1024)
	n, err := fd.Read(readBuffer)
	if err != nil {
		return []string{}, err
	}
	msg := string(readBuffer[0:n])
	cmd := strings.SplitN(msg, " ", 3)
	return cmd, nil
}
func (euw *EtcdUdsWatcher) Find(cmd []string) (string, error) {
	r, ok := euw.rawConfig.Load(cmd[1])
	log.Debug(fmt.Sprintf("%#v,%#v",cmd,r))
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
	case typeJson, typeToml, typeYaml:
		return gjson.Get(raw, cmd[2]).String()
	}
	return raw
}
