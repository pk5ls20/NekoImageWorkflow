package pprof

import (
	"errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/http/pprof"
	"os"
	"time"
)

func RunPprof(addr string) {
	mux := http.NewServeMux()
	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server := http.Server{Addr: addr, Handler: mux}
	logrus.Infof("pprof debug server start!: %v/debug/pprof", addr)
	if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logrus.Error(err)
		logrus.Errorf("Failed to start pprof server, please check if the port %s is occupied.", addr)
		logrus.Error("will exit in 5 seconds...")
		time.Sleep(time.Second * 5)
		os.Exit(1)
	} else {
		logrus.Infof("pprof server started at %v/debug/pprof", addr)
	}
}
