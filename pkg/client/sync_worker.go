package client

import (
	"context"
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
	latestTime   sync.Map
	failConfigs  sync.Map
}

type failConfig struct {
	t time.Time
	k interface{}
	v interface{}
}

func (sw *SyncWorker) SyncOne(key, value interface{}) {
	sk := strings.SplitN(key.(string), "/", 2)
	ext := sk[0] //extension
	rk  := sk[1]  //real key
	memoryFile := "/dev/shm/" + rk
	//mfs := strings.Split(sw.shmfile, "_")
	//modFile := mfs[len(mfs)-1]
	in := sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		//1.Read and write
		fh, err := os.OpenFile(memoryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error("File not open correctly:" + err.Error())
			return
		} else {
			_, err = fh.Write([]byte(value.(string)))
		}
		defer func() {
			err := fh.Close()
			if err != nil {
				log.Error(err.Error())
			}
		}()
		if err == nil {
			//2.Move
			p := time.Now().Format("200601021504")
			err = runCtxCommand("cp", "-f", memoryFile, "/tmp/config_"+rk+"_"+p)
			if err == nil {
				//3.Symlink
				log.Info("ln" + " -sfT" + " /tmp/config_"+rk+"_"+p + " " +sw.storeDir+"/"+rk+"."+ext)
				err = runCtxCommand("ln", "-sfT", "/tmp/config_"+rk+"_"+p, sw.storeDir+"/"+rk+"."+ext)
				if err == nil {
					sw.setLatestTime(rk, time.Now())
				}
			}
		}
		if err != nil {
			log.Error(err.Error())
			sw.setFailConfig(key, value)
		}
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {
		//never invoke..
	}))
	if !in {
		sw.setFailConfig(key, value)
	}
}
func (sw *SyncWorker) SyncAll(configs sync.Map) {
	configs.Range(func(key, value interface{}) bool {
		sw.SyncOne(key, value)
		return true
	})
}

func (sw *SyncWorker) setLatestTime(key string, t time.Time) {
	sw.latestTime.Store(key, t)
}

func (sw *SyncWorker) setFailConfig(key, value interface{}) {
	sw.failConfigs.Store(key, failConfig{
		t: time.Now(),
		k: key,
		v: value,
	})
}

func (sw *SyncWorker) retryFails() {
	tm := timer.NewTimer()
	tm.Ready()
	tm.SetInterval(sw.retrySeconds, func() {
		sw.failConfigs.Range(func(key, value interface{}) bool {
			if lt, ok := sw.latestTime.Load(key); ok {
				vf := value.(failConfig)
				if vf.t.After(lt.(time.Time)) {
					log.Warn(fmt.Sprintf("Try to save configs fail to record last time : %v", vf.k))
					sw.SyncOne(vf.k, vf.v)
				}
			}
			return true
		})
	})
}

func runCtxCommand(commands ...string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	cmd := exec.CommandContext(ctx, commands[0], commands[1:]...)
	err := cmd.Run()
	cancel()
	return err
}
