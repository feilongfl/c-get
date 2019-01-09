package source

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/feilongfl/c-get/core"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var ParseKANMANHUA = Parse_s{
	name:               "看漫画",
	Id:                 0,
	regex:              []string{"tw.manhuagui.com", "www.manhuagui.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoKanmanhua,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterKanmanhua,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImageKanmanhua,
	getImage:           getImageDefault,
}

func getComicInfoKanmanhua(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	comicInfo = comicInfo_s{
		title:    doc.Find("div.book-title > h1").Text(),
		coverUrl: doc.Find("div.book-cover.fl > p > img").AttrOr("src", _unknowPic_),
		info:     doc.Find(".intro-cut").Text(),
		//author:         doc.Find("tr:nth-child(3) > td > a").Text(),
		//classification: doc.Find("tr:nth-child(7) > td > a").Text(),
		//isFinish:       doc.Find("tr:nth-child(5) > td > a").Text() != "连载中",
		//tag:            doc.Find("div.line_height_content").Text(),
	}

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterKanmanhua(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	comicChapter = make([]comicChapter_s, 0)
	//docS := doc.Find("div.zj_list_con").First()
	doc.Find("#chapter-list-0 > ul > li").Each(func(i int, selection *goquery.Selection) {
		adoc := selection.Find("a")
		var c = comicChapter_s{
			name:  adoc.AttrOr("title", ""),
			url:   "https://tw.manhuagui.com" + adoc.AttrOr("href", ""),
			group: 0,
		}
		if c.url != "" {
			comicChapter = append(comicChapter, c)
		}
	})

	log.WithFields(log.Fields{
		"comicChapter length": len(comicChapter),
		"comicChapter":        comicChapter,
	}).Debug("info-comic")

	return comicChapter, nil
}

func KanmanhuaDecode(input string) (result string, err error) {
	re := regexp.MustCompile(`\d,('(.*?)'\[.*?\))`)
	lzdata := re.FindStringSubmatch(input)
	if len(lzdata) != 3 {
		return result, errors.New("regex not match")
	}
	result, err = core.LzDecompressFromBASE64(lzdata[2])

	result = strings.Replace(input, lzdata[1],
		fmt.Sprintf(`'%s'.split('|')`, result), 1)

	result = strings.Replace(result, `window["\x65\x76\x61\x6c"]`, "eval", 1)

	return result, err
}

func getChapterImageKanmanhua(doc *goquery.Document) (imageUrl []string, err error) {
	if doc == nil {
		return nil, errors.New("doc is nil")
	}
	re := regexp.MustCompile(`window.*\)`)
	picJs := re.FindString(doc.Text())
	if picJs == "" {
		return imageUrl, errors.New("regex not match")
	}

	picJs, err = KanmanhuaDecode(picJs)
	picJson, err := core.EvalDecodeNew(picJs)

	re = regexp.MustCompile(`.*?({.*}).*`)
	picJsons := re.FindStringSubmatch(picJson)
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match")
	}
	picJson = picJsons[1]

	var f interface{}
	err = json.Unmarshal([]byte(picJson), &f)
	m := f.(map[string]interface{})

	files := m["files"].([]interface{})
	for _, u := range files {
		image := fmt.Sprintf("https://i.hamreus.com%s%s?cid=%d&md5=%s",
			m["path"],                 //path
			u.(string),                //file
			int((m["cid"]).(float64)), //cid
			(m["sl"].(map[string]interface{}))["md5"], //md5
		)
		imageUrl = append(imageUrl, image)
	}
	return imageUrl, nil
}
