package main

import (
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/dmfed/quickstub"
	"gopkg.in/yaml.v3"
)

//go:embed sample_config.yaml
var sampleConfig string

func main() {
	var (
		confFile   string
		showSample bool
	)
	flag.StringVar(&confFile, "conf", "quickstub.yaml", "yaml file with stubs configuration")
	flag.BoolVar(&showSample, "sample", false, "print sample config and exit")
	flag.Parse()

	if showSample {
		fmt.Println(sampleConfig)
		return
	}

	b, err := os.ReadFile(confFile)
	if err != nil {
		log.Println("error reading config file", confFile, err)
		return
	}

	log.Println("using config file:", confFile)

	conf := new(quickstub.Config)
	if err := yaml.Unmarshal(b, conf); err != nil {
		log.Println("error parsing config:", err)
		return
	}

	srv, err := quickstub.NewServer(conf)
	if err != nil {
		log.Println("error creating new server:", err)
		return
	}
	log.Println("starting quickstub server on", conf.ListenAddr)

	log.Println(srv.ListenAndServe())
}
