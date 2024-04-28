package log

import (
	"github.com/cloudwego/kitex/pkg/klog"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/sirupsen/logrus"
	"os"
)

func init() {
	// init logrus
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&CustomFormatter{})
	// init klog
	logger := kitexlogrus.NewLogger()
	logger.SetOutput(os.Stdout)
	logger.SetLevel(klog.LevelDebug)
	logger.Logger().SetReportCaller(true)
	logger.Logger().SetFormatter(&CustomFormatter{})
	klog.SetLogger(logger)
}
