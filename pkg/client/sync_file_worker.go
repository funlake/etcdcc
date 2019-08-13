package client

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/funlake/etcdcc/pkg/log"
	"github.com/funlake/gopkg/jobworker"
	"github.com/funlake/gopkg/timer"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"
)

//SyncFileWorker : Worker of sync configuration to file
type SyncFileWorker struct {
	storeDir     string
	shmfile      string
	timeout      int
	retrySeconds int
	dispatcher   *jobworker.BlockingDispatcher
	latestTime   sync.Map
	failQueue    sync.Map
}

//RemoveOne : Remove one configuration setting
func (sw *SyncFileWorker) RemoveOne(key interface{}) {
	ext, rk := getKeyAndExt(key.(string))
	err := os.Remove(sw.storeDir + "/" + rk + "." + ext)
	if err != nil {
		log.Error(err.Error())
	}

}

//SyncOne : Sync one configuration setting
func (sw *SyncFileWorker) SyncOne(key, value interface{}) {
	ext, rk := getKeyAndExt(key.(string))
	rk = strings.Replace(rk, "/", "_", -1)
	memoryFile := "/dev/shm/" + rk
	in := sw.dispatcher.Put(jobworker.NewSimpleJob(func() {
		//1.Read and write
		fh, err := os.OpenFile(memoryFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			log.Error("File not open correctly:" + err.Error())
			return
		}
		vb, _ := base64.StdEncoding.DecodeString(value.(string))
		_, err = fh.Write(vb)
		defer func() {
			err := fh.Close()
			if err != nil {
				log.Error(err.Error())
			}
		}()
		if err == nil {
			//2.Move
			p := time.Now().Format("200601021504")
			err = sw.moveFile(memoryFile, "/tmp/config_"+rk+"_"+p)
			if err == nil {
				//3.Symlink
				err = sw.linkFile("/tmp/config_"+rk+"_"+p, sw.storeDir+"/"+rk+"."+ext)
				if err == nil {
					sw.setLatestTime(rk, time.Now())
				}
			}
		}
		if err != nil {
			log.Error(err.Error())
			sw.pushToFailQueue(key, value)
		}
	}, func() string {
		return "config_sync_job"
	}, func(dispatcher *jobworker.BlockingDispatcher) {
		//never invoke..
	}))
	if !in {
		sw.pushToFailQueue(key, value)
	}
}
func (sw *SyncFileWorker) setLatestTime(key string, t time.Time) {
	sw.latestTime.Store(key, t)
}

func (sw *SyncFileWorker) pushToFailQueue(key, value interface{}) {
	sw.failQueue.Store(key, failConfig{
		t: time.Now(),
		k: key,
		v: value,
	})
}

func (sw *SyncFileWorker) moveFile(source, target string) error {
	return runCtxCommand("cp", "-f", source, target)
}

func (sw *SyncFileWorker) linkFile(source, target string) error {
	return runCtxCommand("ln", "-sfT", source, target)
}

//Retry implements
func (sw *SyncFileWorker) Retry() {
	tm := timer.NewTimer()
	tm.Ready()
	tm.SetInterval(sw.retrySeconds, func() {
		sw.failQueue.Range(func(key, value interface{}) bool {
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

func getKeyAndExt(key string) (string, string) {
	sk := strings.SplitN(key, "/", 2)
	ext := sk[0] //extension
	rk := sk[1]  //real key
	return ext, rk
}
