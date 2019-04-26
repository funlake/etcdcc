package log
import(
	"fmt"
	"github.com/rs/zerolog"
	"os"
	"strings"
)
const LOGTIMEFORMAT  = "2006-01-02 15:04:05"
var log zerolog.Logger
func init()  {
	zerolog.CallerSkipFrameCount = 3
	output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: LOGTIMEFORMAT}
	output.FormatLevel = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf(" | %s", i))
	}
	output.FormatMessage = func(i interface{}) string {
		return fmt.Sprintf(" | %s", i)
	}
	output.FormatFieldName = func(i interface{}) string {
		return fmt.Sprintf(" %s:", i)
	}
	output.FormatFieldValue = func(i interface{}) string {
		return strings.ToUpper(fmt.Sprintf("%s ", i))
	}
	output.FormatCaller = func(i interface{}) string {
		var c string
		if cc, ok := i.(string); ok {
			c = cc
		}
		if len(c) > 0 {
			cwd, err := os.Getwd()
			if err == nil {
				c = strings.TrimPrefix(c, cwd)
				c = strings.TrimPrefix(c, "/")
			}
		}
		return "| "+c
	}
	log = zerolog.New(output).With().Timestamp().Logger()

}
func Debug(msg string,fields ... map[string]interface{}){
	log.Debug().Fields(fields[0]).Caller().Msg(msg)
}
func Info(msg string){
	log.Info().Caller().Msg(msg)
}
func Warn(msg string){
	log.Warn().Caller().Msg(msg)
}
func Error(msg string){
	log.Error().Caller().Msg(msg)
}
func Fatal(msg string){
	log.Fatal().Caller().Msg(msg)
}