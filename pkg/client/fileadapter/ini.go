package fileadapter

import "sync"

type Ini struct {
	filename string
}

func (ini *Ini) SetFileName(filename string) {
	ini.filename = filename
}
func (ini *Ini) Save(configs sync.Map) {

}
