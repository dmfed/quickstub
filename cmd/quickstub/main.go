package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/dmfed/quickstub"
)

//go:embed sample_config.yaml
var sampleConfig string

//go:embed help.txt
var helpMessage string

func main() {
	var (
		confFile   string
		showSample bool
		certFile   string
		pkeyFile   string
	)
	flag.StringVar(&confFile, "conf", "", "yaml file with stubs configuration")
	flag.StringVar(&certFile, "cert", "", "RSA certificate file to use for TLS")
	flag.StringVar(&pkeyFile, "pkey", "", "RSA private key file to use for TLS")
	flag.BoolVar(&showSample, "sample", false, "print sample config and exit")
	flag.Usage = func() {
		fmt.Println(helpMessage)
		flag.PrintDefaults()
	}
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

	conf, err := quickstub.ParseConfig(b)
	if err != nil {
		log.Println("error parsing config:", err)
		return
	}

	app, err := quickstub.NewQuickStub(conf)
	if err != nil {
		log.Println("error creating new server:", err)
		return
	}

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
		<-c
		ctx, stop := context.WithTimeout(context.Background(), time.Second*15)
		defer stop()
		app.Shutdown(ctx)
	}()

	log.Println("starting quickstub server on", conf.ListenAddr)
	log.Println(app.ListenAndServe())
}
