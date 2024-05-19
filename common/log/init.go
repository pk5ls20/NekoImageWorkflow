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
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&CustomFormatter{})
	// init klog
	logger := kitexlogrus.NewLogger()
	logger.SetOutput(os.Stdout)
	logger.Logger().SetReportCaller(true)
	logger.Logger().SetFormatter(&CustomFormatter{})
	klog.SetLogger(logger)
	// Read log level from environment variable
	logLevel := os.Getenv("NEKOIMAGEWORKFLOW_LOGLEVEL")
	switch logLevel {
	case "TRACE":
		logrus.SetLevel(logrus.TraceLevel)
		logger.SetLevel(klog.LevelTrace)
	case "DEBUG":
		logrus.SetLevel(logrus.DebugLevel)
		logger.SetLevel(klog.LevelDebug)
	case "INFO":
		logrus.SetLevel(logrus.InfoLevel)
		logger.SetLevel(klog.LevelInfo)
	case "WARN":
		logrus.SetLevel(logrus.WarnLevel)
		logger.SetLevel(klog.LevelWarn)
	case "ERROR":
		logrus.SetLevel(logrus.ErrorLevel)
		logger.SetLevel(klog.LevelError)
	case "FATAL":
		logrus.SetLevel(logrus.FatalLevel)
		logger.SetLevel(klog.LevelFatal)
	default:
		logrus.SetLevel(logrus.DebugLevel)
		logger.SetLevel(klog.LevelDebug)
	}
}
