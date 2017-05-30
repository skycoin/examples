package main

import (
	"github.com/skycoin/skycoin/src/util"
	"log"
	"os"
	"os/signal"
)

func main() {
	quit := CatchInterrupt()
	config, e := NewConfig().Parse().PostProcess()
	if e != nil {
		panic(e)
	}
	util.InitDataDir(config.ConfigDir())
	defer log.Println("Goodbye.")

	log.Println("[CONFIG] Starting cxo client and server on port", config.CXOPort())
	container, e := NewContainer(config)
	CatchError(e, "unable to create cxo container")
	defer container.Close()

	log.Println("!!! EVERYTHING UP AND RUNNING !!!")
	defer log.Println("Shutting down...")
	<-quit
}

// CatchInterrupt catches Ctrl+C behaviour.
func CatchInterrupt() chan int {
	quit := make(chan int)
	go func(q chan<- int) {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, os.Interrupt)
		<-sigchan
		signal.Stop(sigchan)
		q <- 1
	}(quit)
	return quit
}

// CatchError catches an error and panics.
func CatchError(e error, msg string, args ...interface{}) {
	if e != nil {
		log.Panicf(msg+": %v", append(args, e)...)
	}
}
