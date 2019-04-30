package client

import (
	"context"
	"etcdcc/apiserver/pkg/client/fileadapter"
	"etcdcc/apiserver/pkg/log"
	"fmt"
	"github.com/funlake/gopkg/jobworker"
	"github.com/funlake/gopkg/timer"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

type SyncWorker struct {
	storeDir     string
	shmfile      string
	timeout      int
	retrySeconds int
	dispatcher   *jobworker.BlockingDispatcher
	latestTime   time.Time
	*failConfig
}

type failConfig struct {
	t time.Time
	c sync.Map
}

func (sw *SyncWorker) Do(configs sync.Map) {
	memoryFile := "/dev/shm/" + sw.shmfile
	mfs := strings.Split(sw.shmfile, "_")
	modFile := mfs[len(mfs)-1]
	in := sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		//1.Read and write
		fh, err := os.OpenFile(memoryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Fatal("File not open correctly:" + err.Error())
		}
		fa := fileadapter.Json{File: fh}
		err = fa.Save(configs)
		defer func() {
			err := fh.Close()
			if err != nil {
				log.Error(err.Error())
			}
		}()
		if err == nil {
			//2.Move
			p := time.Now().Format("200601021504")
			err = runCtxCommand("cp", "-f", memoryFile, "/tmp/config_"+p)
			if err == nil {
				//3.Symlink
				err = runCtxCommand("ln", "-sfT", "/tmp/config_"+p, sw.storeDir+"/"+modFile+".json")
				if err == nil {
					sw.setLatestTime(time.Now())
				}
			}
		}
		if err != nil {
			log.Error(err.Error())
			sw.setFailConfig(configs)
		}
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {
		//never invoke..
	}))
	if !in {
		sw.setFailConfig(configs)
	}
}

func (sw *SyncWorker) setLatestTime(t time.Time) {
	sw.latestTime = t
}

func (sw *SyncWorker) setFailConfig(configs sync.Map) {
	sw.failConfig = &failConfig{
		t: time.Now(),
		c: configs,
	}
}

func (sw *SyncWorker) retryFails() {
	tm := timer.NewTimer()
	tm.Ready()
	tm.SetInterval(sw.retrySeconds, func() {
		if sw.failConfig != nil {
			if sw.failConfig.t.After(sw.latestTime) {
				log.Warn(fmt.Sprintf("Try to save configs fail to record last time : %v", sw.failConfig.c))
				sw.Do(sw.failConfig.c)
			} else {
				//discard,since configs already the newer one
				sw.failConfig = nil
			}
		}
	})
}

func runCtxCommand(commands ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	err := cmd.Run()
	cancel()
	return err
}
