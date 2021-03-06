package source

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/feilongfl/c-get/core"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var ParseTENCENT = Parse_s{
	name:               "腾讯漫画",
	Id:                 0,
	regex:              []string{"ac.qq.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoTencent,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterPufei,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImagePufei,
	getImage:           getImageDefault,
}

func getComicInfoTencent(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	dinfo := doc.Find("div.works-intro")
	comicInfo = comicInfo_s{
		coverUrl: dinfo.Find("div.works-cover > a > img").AttrOr("src", _unknowPic_),
	}
	dinfotext := dinfo.Find("div.works-intro-text")
	comicInfo.title = dinfotext.Find("div.works-intro-title > strong").Text()

	dinfotext = dinfo.Find("div.works-intro-text")
	comicInfo.author = dinfotext.Find("p.works-intro-digi > span > em").Text()

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterTencent(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	comicChapter = make([]comicChapter_s, 0)
	docS := doc.Find("div.zj_list_con").First()
	docS.Find("ul > li").Each(func(i int, selection *goquery.Selection) {
		adoc := selection.Find("a")
		var c = comicChapter_s{
			name:  adoc.AttrOr("title", ""),
			url:   adoc.AttrOr("href", ""),
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

func getChapterImageTencent(doc *goquery.Document) (imageUrl []string, err error) {
	//fmt.Print(doc.Text())
	re := regexp.MustCompile("eval\\(.*\\)")
	picJs := re.FindString(doc.Text())
	//log.Info(picJs)
	picJson, err := core.EvalDecode(picJs)
	picJson, err = core.EvalDecode(picJson)
	re = regexp.MustCompile("var pages='(.*)'")
	picJsons := re.FindStringSubmatch(picJson)
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match!")
	}
	picJson = picJsons[1]
	var f map[string]interface{}
	err = json.Unmarshal([]byte(picJson), &f)
	imageUrl = strings.Split(f["page_url"].(string), "\r\n")
	for index, url := range imageUrl {
		imageUrl[index] = "https://images.Tencent.com/" + url
	}
	return imageUrl, nil
}
