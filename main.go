package main

import (
	"c-get/core"
	"c-get/source"

	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

func initParse() {
	parseArray := []source.Parse_s{
		// add parse here
		source.ParseDemo,
		source.ParseDMZJ,
		source.ParsePUFEI,
		source.ParseTENCENT,
		source.ParseDMZJV2,
		source.ParseKANMANHUA,
	}

	for i, p := range parseArray {
		p.Id = i
		source.ParseList = append(source.ParseList, p)
	}
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)
	log.SetLevel(log.WarnLevel)
	//log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
	initParse()
	core.CliRecv()
}
