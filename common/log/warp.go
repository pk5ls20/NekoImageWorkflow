package log

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
)

func ErrorWrap(err error) error {
	if err != nil {
		_, file, line, ok := runtime.Caller(1)
		var logMsg string
		if ok {
			logMsg = fmt.Sprintf("<ERRORWARP> Error %s raised in %s:%d", err.Error(), file, line)
		} else {
			logMsg = "<ERRORWARP> Could not get caller information"
		}
		logrus.Debug(logMsg)
	}
	return err
}
