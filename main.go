package main

import (
	"github.com/feilongfl/c-get/command"
	"github.com/feilongfl/c-get/source"

	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

var parseArray = []source.Parse_s{
	// add parse here
	source.ParseDemo,
	source.ParseDMZJ,
	source.ParsePUFEI,
	source.ParseTENCENT,
	source.ParseDMZJV2,
	source.ParseKANMANHUA,
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	//log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
	source.InitParse(parseArray)
	command.CliRecv()
}
