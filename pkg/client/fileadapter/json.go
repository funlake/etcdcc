package fileadapter

import (
	"etcdcc/apiserver/pkg/log"
	"os"

	//"encoding/json"
	"github.com/json-iterator/go"
	"sync"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

type Json struct {
	File *os.File
}

func (j *Json) SetFileHandler(file *os.File) {
	j.File = file
}
func (j *Json) Save(configs sync.Map) error {
	cfg := make(map[interface{}]interface{})
	configs.Range(func(key, value interface{}) bool {
		cfg[key] = value
		return true
	})
	c, err := json.Marshal(cfg)
	if err != nil {
		log.Error(err.Error())
		return err
	} else {
		_, err := j.File.Write(c)
		if err != nil {
			log.Error(err.Error())
			return err
		}
	}
	return nil
}
