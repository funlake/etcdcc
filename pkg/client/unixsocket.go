package client

import (
	"bytes"
	"encoding/json"
	"github.com/BurntSushi/toml"
	"github.com/funlake/etcdcc/pkg/log"
	"github.com/ghodss/yaml"
	"net"
	"os"
	"strings"
)

type UnixSocket struct {
	Wch *MemoryWatcher
}

func (uxs *UnixSocket) Serve(sockFile string) {
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
		go uxs.handle(fd)
	}
}
func (uxs *UnixSocket) handle(fd net.Conn) {
	defer fd.Close()
	for {
		cmd, err := uxs.getCmd(fd)
		if err != nil {
			return
		}
		switch strings.ToLower(cmd[0]) {
		case "get":
			val, err := uxs.Wch.Find(cmd)
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
			r, ok := uxs.Wch.rawConfig.Load(cmd[1])
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
func (uxs *UnixSocket) getCmd(fd net.Conn) ([]string, error) {
	readBuffer := make([]byte, 1024)
	n, err := fd.Read(readBuffer)
	if err != nil {
		return []string{}, err
	}
	msg := string(readBuffer[0:n])
	cmd := strings.SplitN(msg, " ", 3)
	return cmd, nil
}
