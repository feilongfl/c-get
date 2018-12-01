package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
)

type parse_s struct {
	name               string
	id                 int
	regex              []string
	getComicInfoReq    func(url string) (*goquery.Document, error)
	getComicInfo       func(doc *goquery.Document) (comicInfo comicInfo_s, err error)
	getComicChapterReq func(url string) (*goquery.Document, error)
	getComicChapter    func(doc *goquery.Document) (comicChapter []comicChapter_s, err error)
	getChapterImageReq func(url string, referer string) (*goquery.Document, error)
	getChapterImage    func(doc *goquery.Document) (imageUrl []string, err error)
	getImage           func(imageUrl string, referer string, path string) (err error)
}

var parseList = make([]parse_s, 0)

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

func newParseFromUrl(url string) (p *parse_s, err error) {
	for _, p := range parseList {
		for _, r := range p.regex {
			m, err := regexp.MatchString(r, url)
			if err != nil {
				log.Fatal(err)
			}
			if m {
				return &p, nil
			}
		}
	}
	return nil, errors.New("no parse match")
}

//处理信息页网络下载请求
func getComicInfoReqDefault(url string) (*goquery.Document, error) {
	var _httpReq = defaultHttpReq
	_httpReq.url = url
	doc, err := htmlGet(_httpReq)
	return doc, err
}

func getChapterImageReqDefault(url string, referer string) (*goquery.Document, error) {
	var _httpReq = defaultHttpReq
	_httpReq.url = url
	_httpReq.referer = referer
	doc, err := htmlGet(_httpReq)
	return doc, err
}

//解析信息页
func getComicInfoDefault(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	return comicInfo, errors.New("nil getComicInfoDefault func")
}

//解析章节列表
func getComicChapterDefault(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	return comicChapter, errors.New("nil getComicChapterDefault func")
}

//解析图片列表
func getChapterImageDefault(doc *goquery.Document) (imageUrl []string, err error) {
	return imageUrl, errors.New("nil getChapterImageDefault func")
}

func getImageDefault(imageUrl string, referer string, path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// path/to/whatever does not exist
		var _httpReq = defaultHttpReq
		_httpReq.url = "https://images.dmzj.com/" + imageUrl
		_httpReq.referer = referer
		err = fileGet(_httpReq, path)
		if err != nil {
			os.Remove(path)
		}
		return err
	} else {
		return nil
	}
}
