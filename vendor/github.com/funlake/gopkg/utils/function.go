package utils

import (
	"bytes"
	"encoding/json"
	"github.com/funlake/gopkg/utils/log"
	"unsafe"
)

func RoutineRecover(msg string) {
	if err := recover(); err != nil {
		log.Error("%s : %s ", msg, err)
	}
}
func WrapGo(fun func(), msg string) {
	go func() {
		defer RoutineRecover(msg)
		fun()
	}()
}
func StrToByte(s string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&s))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func ByteToStr(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

//设置html过滤为false,否则默认的json marshal会过滤html字符,
//导致返回数据不是raw string
func JSONMarshal(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}
