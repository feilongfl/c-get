package main

import (
	//"log"
	log "github.com/sirupsen/logrus"
	"os"
)

var _version_ = "v1.0.0"
var _commit_ = "manual"
var _unknowPic_ = "http://xxx.pic"

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

	id += 1
	parseList = append(parseList, parse_s{
		name:               "扑飞漫画",
		id:                 id,
		regex:              []string{"www.pufei.net"},
		getComicInfoReq:    getComicInfoReqDefault,
		getComicInfo:       getComicInfoPufei,
		getComicChapterReq: getComicInfoReqDefault,
		getComicChapter:    getComicChapterPufei,
		getChapterImageReq: getChapterImageReqDefault,
		getChapterImage:    getChapterImagePufei,
		getImage:           getImageDefault,
	})

	id += 1
	parseList = append(parseList, parse_s{
		name:               "t腾讯漫画",
		id:                 id,
		regex:              []string{"ac.qq.com"},
		getComicInfoReq:    getComicInfoReqDefault,
		getComicInfo:       getComicInfoTencent,
		getComicChapterReq: getComicInfoReqDefault,
		getComicChapter:    getComicChapterPufei,
		getChapterImageReq: getChapterImageReqDefault,
		getChapterImage:    getChapterImagePufei,
		getImage:           getImageDefault,
	})

	id += 1
	parseList = append(parseList, parse_s{
		name:               "动漫之家v2",
		id:                 id,
		regex:              []string{"manhua.dmzj.com"},
		getComicInfoReq:    getComicInfoReqDefault,
		getComicInfo:       getComicInfoDmzjV2,
		getComicChapterReq: getComicInfoReqDefault,
		getComicChapter:    getComicChapterDmzjV2,
		getChapterImageReq: getChapterImageReqDefault,
		getChapterImage:    getChapterImageDmzjV2,
		getImage:           getImageDefault,
	})
}

func main() {
	//log.SetFormatter(&log.JSONFormatter{})
	//runtime.GOMAXPROCS(runtime.NumCPU())
	log.SetOutput(os.Stdout)
	//log.SetLevel(log.WarnLevel)
	log.SetLevel(log.InfoLevel)
	//log.SetLevel(log.DebugLevel)
	initParse()
	cliRecv()
}
