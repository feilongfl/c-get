package main

import (
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

var _version_ = "v1.0.0"
var _commit_ = "manual"
var _unknowPic_ = "http://xxx.pic"

var ParseDemo = parse_s{
	name:               "demo",
	id:                 0,
	regex:              []string{"www.example.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoDefault,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterDefault,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImageDefault,
	getImage:           getImageDefault,
}

func initParse() {
	parseArray := []parse_s{
		// add parse here
		ParseDemo,
		ParseDMZJ,
		ParsePUFEI,
		ParseTENCENT,
		ParseDMZJV2,
		ParseKANMANHUA,
	}

	for i, p := range parseArray {
		p.id = i
		parseList = append(parseList, p)
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
	cliRecv()
}
