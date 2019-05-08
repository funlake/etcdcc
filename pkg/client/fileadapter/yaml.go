package fileadapter

import (
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"sync"
)

type Yaml struct {
	File *os.File
}

func (y *Yaml) SetFileHandler(file *os.File) {
	y.File = file
}
func (y *Yaml) Save(configs sync.Map) error {
	//cnt	:= ``
	//configs.Range(func(key, value interface{}) bool {
	//	cnt = cnt + fmt.Sprintf("%s: %s\n",key,value)
	//	return true
	//})
	out, err := yaml.Marshal(configs)
	if err != nil{
		log.Error(err.Error())
		return err
	}
	log.Info(fmt.Sprintf("%x",out))
	_,err = y.File.Write(out)
	if err != nil{
		log.Error(err.Error())
		return err
	}
	return nil
}