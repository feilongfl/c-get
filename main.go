package main

import (
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
)

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
	initParse()
	cliRecv()
}
