package client

import (
	"context"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/funlake/gopkg/jobworker"
	"os"
	"os/exec"
	"sync"
	"time"
)

type SyncWorker struct {
	second int
	dispatcher *jobworker.BlockingDispatcher
}

func (sw *SyncWorker) Do (configs sync.Map){
	in := sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		//1.Read and write
		fh,err := os.OpenFile("/dev/shm/config",os.O_RDWR | os.O_CREATE,0666)
		if err != nil {
			log.Fatal("File not open correctly:"+err.Error())
		}
		configs.Range(func(key, value interface{}) bool {
			log.Info(fmt.Sprintf("%s:%s",key,value))
			_,_ = fh.Write([]byte(fmt.Sprintf("%s=%s\n",key,value)))
			return true
		})
		defer func() {
			err := fh.Close()
			if err != nil{
				log.Error(err.Error())
			}
		}()
		if err == nil {
			//2.Move
			p := time.Now().Format("200601021504")
			ctx,cancel := context.WithTimeout(context.Background(),3 * time.Second)
			cmd := exec.CommandContext(ctx,"cp", "-f", "/dev/shm/config", "/tmp/config_"+p)
			err = cmd.Run()
			cancel()
			if err != nil {
				log.Error(err.Error())
			}
			if err == nil {
				//3.Symlink
				ctx,cancel := context.WithTimeout(context.Background(),3 * time.Second)
				cmd := exec.CommandContext(ctx,"ln", "-sfT", "/tmp/config_"+p, "/opt/test.ini")
				err = cmd.Run()
				cancel()
				if err != nil {
					log.Error(err.Error())
				}
			}
		}
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {
		//never invoke..
	}))
	if !in{
		//todo : failure handler
	}
}