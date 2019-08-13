package client

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/funlake/etcdcc/pkg/log"
	"github.com/funlake/gopkg/cache"
	"github.com/ghodss/yaml"
	"github.com/tidwall/gjson"
	"github.com/zieckey/goini"
	"strings"
	"sync"
)

//EtcdUdsWatcher : Unix domain socket watcher for etcd
type EtcdMemoryWatcher struct {
	GeneralWatcher
	rawConfig sync.Map
	Tc        *cache.TimerCacheEtcd
}

//KeepEyesOnKey : Watching specific key
func (emw *EtcdMemoryWatcher) KeepEyesOnKey(key string) {}

//KeepEyesOnKeyWithPrefix : Watch etcd with prefix
func (emw *EtcdMemoryWatcher) KeepEyesOnKeyWithPrefix(prefix string) {
	emw.Init(emw.Tc, prefix, func(k, v string) {
		emw.saveLocal(k, v)
	})
	emw.Watch(emw.Tc, prefix, func(k, v string) {
		emw.saveLocal(k, v)
	}, func(mk, k string, cancel context.CancelFunc) {
		if mk == k {
			cancel()
		}
		emw.rawConfig.Delete(k)
	})
}
func (emw *EtcdMemoryWatcher) saveLocal(k, v string) {
	r, err := base64.StdEncoding.DecodeString(v)
	if err != nil {
		log.Error(fmt.Sprintf("Base64 decode error with key %s : %s:", k, err.Error()))
	} else {
		//transform yaml/toml to json for easy handling
		r, err = emw.jsonEncode(r, k)
		log.Debug(fmt.Sprintf("saving data %s -> %s", k, string(r)))
		emw.rawConfig.Store(k, string(r))
	}
}
func (emw *EtcdMemoryWatcher) jsonEncode(r []byte, prefix string) ([]byte, error) {
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
func (emw *EtcdMemoryWatcher) Find(cmd []string) (string, error) {
	r, ok := emw.rawConfig.Load(cmd[1])
	log.Debug(fmt.Sprintf("%#v,%#v", cmd, r))
	if ok {
		if len(cmd) > 2 {
			val := emw.getSpecifyKey(r.(string), cmd)
			if val != "" {
				return val, nil
			}
		} else {
			return r.(string), nil
		}
	}
	return "", errors.New("NotFound")
}
func (emw *EtcdMemoryWatcher) getSpecifyKey(raw string, cmd []string) string {
	t := strings.SplitN(cmd[1], "/", 2)
	switch t[0] {
	case typeJson, typeToml, typeYaml:
		return gjson.Get(raw, cmd[2]).String()
	}
	return raw
}
