package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"
	"time"
)

type CustomFormatter struct{}

func (f *CustomFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var levelColor int
	var fileLine string
	var output string
	levelText := strings.ToUpper(entry.Level.String())
	if entry.Caller != nil {
		lastPath := strings.Split(entry.Caller.File, "/")
		fileLine = fmt.Sprintf("%s:%d", lastPath[len(lastPath)-1], entry.Caller.Line)
	} else {
		fileLine = "<unknown.go>:0"
	}
	switch entry.Level {
	case logrus.DebugLevel:
		levelColor = 36 // Cyan
	case logrus.InfoLevel:
		levelColor = 32 // Green
	case logrus.WarnLevel:
		levelColor = 33 // Yellow
	case logrus.ErrorLevel:
		levelColor = 31 // Red
	case logrus.FatalLevel, logrus.PanicLevel:
		levelColor = 35 // Magenta
	default:
		levelColor = 37 // White
	}
	if entry.Level == logrus.ErrorLevel || entry.Level == logrus.FatalLevel || entry.Level == logrus.PanicLevel {
		output = fmt.Sprintf("\033[32m%s \033[0m[\033[%dm%s\033[0m]"+
			" \033[34;4m%s\033[0m | \033[38;2;255;0;0;48;2;255;188;212m%s\033[0m\n",
			time.Now().Format("01-02 15:04:05.00000"), levelColor, levelText, fileLine, entry.Message)
	} else {
		output = fmt.Sprintf("\033[32m%s \033[0m[\033[%dm%s\033[0m] \033[34;4m%s\033[0m | %s\n",
			time.Now().Format("01-02 15:04:05.00000"), levelColor, levelText, fileLine, entry.Message)
	}
	return []byte(output), nil
}
