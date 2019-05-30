package jobworker

import (
//"github.com/funlake/gopkg/utils/log"
)

type Dispatcher interface {
	Put(job WorkerJob) bool
	Run(maxWork int)
	Ready()
	Stop()
	StopStatus() bool
}
