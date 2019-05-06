package fileadapter

import (
	"os"
	"sync"
)

type Adapter interface {
	SetFileHandler(file *os.File)
	//configs save to files
	Save(configs sync.Map) error
}
