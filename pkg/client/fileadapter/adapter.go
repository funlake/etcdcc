package fileadapter

import (
	"os"
	"sync"
)

type Adapter interface {
	SetFileName(file *os.File)
	//configs save to files
	Save(configs sync.Map) error
}
