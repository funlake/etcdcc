package log

import (
	"fmt"
	"path/filepath"
	"runtime"
	"time"
)

//var (
//	logger, _ = zap.NewProduction()
//	suger     = logger.Sugar()
//)
const (
	color_red = uint8(iota + 91)
	color_green
	color_yellow
	color_blue
	color_magenta //洋红

	info = "[I]"
	trac = "[T]"
	erro = "[E]"
	warn = "[W]"
	succ = "[S]"
)

// see complete color rules in document in https://en.wikipedia.org/wiki/ANSI_escape_code#cite_note-ecma48-13
func Trace(format string, a ...interface{}) {
	prefix := yellow(trac)
	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
}

func Info(format string, a ...interface{}) {
	prefix := blue(info)
	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	//suger.Infof(format, a...)
	//logger.Sync()
}

func Success(format string, a ...interface{}) {
	prefix := green(succ)
	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	//suger.Infof(format, a...)
	//logger.Sync()
}

func Warning(format string, a ...interface{}) {
	prefix := yellow(warn)
	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	//suger.Warnf( format, a...)
	//logger.Sync()
}

func Error(format string, a ...interface{}) {
	prefix := red(erro)
	fmt.Println(formatLog(prefix), fmt.Sprintf(format, a...))
	//suger.Errorf( format, a...)
	//logger.Sync()
}

func red(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_red, s)
}

func green(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_green, s)
}

func yellow(s string) string {

	return fmt.Sprintf("\x1b[%dm%s\x1b[0m [%s,%s]", color_yellow, s)
}

func blue(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_blue, s)
}

func magenta(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", color_magenta, s)
}

func formatLog(prefix string) string {
	_, fn, line, _ := runtime.Caller(2)
	file := filepath.Base(fn)
	fd := fmt.Sprintf("[%s:%d]", file, line)
	return time.Now().Format("2006/01/02 15:04:05") + " " + prefix + " " + fd
}
