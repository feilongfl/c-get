package main

import (
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

var _version_ = "v1.0.0"
var _commit_ = "manual"

func initParse() {
	parseList = append(parseList, parse_s{
		name:               "demo",
		id:                 -1,
		regex:              []string{"www.example.com"},
		getComicInfoReq:    getComicInfoReqDefault,
		getComicInfo:       getComicInfoDefault,
		getComicChapterReq: getComicInfoReqDefault,
		getComicChapter:    getComicChapterDefault,
		getChapterImageReq: getChapterImageReqDefault,
		getChapterImage:    getChapterImageDefault,
		getImage:           getImageDefault,
	})

	id := 0
	parseList = append(parseList, parse_s{
		name:               "动漫之家",
		id:                 id,
		regex:              []string{"www.dmzj.com"},
		getComicInfoReq:    getComicInfoReqDefault,
		getComicInfo:       getComicInfoDmzj,
		getComicChapterReq: getComicInfoReqDefault,
		getComicChapter:    getComicChapterDmzj,
		getChapterImageReq: getChapterImageReqDefault,
		getChapterImage:    getChapterImageDmzj,
		getImage:           getImageDefault,
	})

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
