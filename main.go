package main

import (
	"flag"
	"net/http"
	"os"
	"runtime/pprof"

	"github.com/opentarock/http-relay/relay"
	"github.com/opentarock/http-relay/vars"
	"github.com/opentarock/service-api/go/log"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	logger := log.New("name", vars.ModuleName)

	// profiliing related flag
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Error("Error creating cpuprofile file", "error", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	logger.Info("Starting http relay ...")

	http.HandleFunc("/relay", relay.RelayHandler)

	err := http.ListenAndServe(":8081", nil)
	if err != nil {
		logger.Error("Error starting http server", "error", err)
		os.Exit(1)
	}

}
