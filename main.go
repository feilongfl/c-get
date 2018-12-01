package main

import (
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

var _version_ = "v1.0.0"
var _commit_ = "manual"

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	//log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
	initParse()
	cliRecv()
}
