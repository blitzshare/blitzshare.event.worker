package main

import (
	"os"
	"os/signal"
	"syscall"

	"blitzshare.event.worker/app"
	"blitzshare.event.worker/app/config"
	"blitzshare.event.worker/app/dependencies"
	log "github.com/sirupsen/logrus"
)

func initLog() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func main() {
	initLog()
	log.Println("Hello from event worker")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config %v\n", err)
	}
	deps, err := dependencies.NewDependencies(cfg)
	if err != nil {
		log.Fatalf("failed to load dependencies %v\n", err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	app.Start(deps)

	log.Printf("worker is running")
	<-signals

	if err != nil {
		log.Fatalf("failed to stop http server %v", err)
	}
}
