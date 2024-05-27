package main

import (
	"flag"
)

var (
	runTime     = flag.Int("runtime", 0, "Run time in seconds, set 0 to disable")
	pprofEnable = flag.Bool("pprof", false, "Enable pprof")
	pprofAddr   = flag.String("pprofAddr", "127.0.0.1:6060", "pprof address")
)

func parseFlags() {
	flag.Parse()
}
