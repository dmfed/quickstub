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
		fmt.Fprint(os.Stdout, sampleConfig)
		return
	}

	b, err := os.ReadFile(confFile)
	if err != nil {
		log.Println(err)
		return
	}

	conf := new(quickstub.Config)
	if err := yaml.Unmarshal(b, conf); err != nil {
		log.Println(err)
		return
	}

	srv, err := quickstub.NewServer(conf)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("starting quickstub server on %s\n", conf.ListenAddr)
	log.Println(srv.ListenAndServe())
}
