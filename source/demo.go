package source

import (
	"c-get/part3rd"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"os"
)

var ParseDemo = Parse_s{
	name:               "demo",
	Id:                 0,
	regex:              []string{"www.example.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoDefault,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterDefault,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImageDefault,
	getImage:           getImageDefault,
}

//处理信息页网络下载请求
func getComicInfoReqDefault(url string) (*goquery.Document, error) {
	var _httpReq = part3rd.DefaultHttpReq(url)
	doc, err := part3rd.HtmlGet(_httpReq)
	return doc, err
}

func getChapterImageReqDefault(url string, referer string) (*goquery.Document, error) {
	var _httpReq = part3rd.DefaultHttpReq(url)
	_httpReq.Referer = referer
	doc, err := part3rd.HtmlGet(_httpReq)
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
		var _httpReq = part3rd.DefaultHttpReq(imageUrl)
		_httpReq.Referer = referer
		err = part3rd.FileGet(_httpReq, path)
		if err != nil {
			os.Remove(path)
		}
		return err
	} else {
		return nil
	}
}
