package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"gopkg.in/yaml.v3"

	"vl3nse/pkg/nseconfig"
)

var (
	configFile = flag.String("cfg", "/etc/cnns-nse/config.yaml", "the path to config file, default: /etc/cnns-nse/config.yaml")

	logger = log.New(os.Stderr, "test", 0)
)

func main() {
	flag.Parse()

	f, err := os.Open(*configFile)
	fatalOnErr(err)

	cfg := &nseconfig.Config{}
	err = nseconfig.NewConfig(yaml.NewDecoder(f), cfg)
	fatalOnErr(err)

	//TODO the connection logic, for now just spew the config
	err = json.NewEncoder(os.Stdout).Encode(cfg)
	fatalOnErr(err)

	<-newOSSignalChannel()
}

func fatalOnErr(err error) {
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func newOSSignalChannel() chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c,
		os.Interrupt,
		// More Linux signals here
		syscall.SIGHUP,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	return c
}
