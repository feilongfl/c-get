package main

import (
	"encoding/json"
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

func getComicInfoDmzj(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	dinfo := doc.Find("div.wrap_intro_l")
	comicInfo = comicInfo_s{
		title:    dinfo.Find("div.comic_deCon > h1").Text(),
		coverUrl: dinfo.Find("div.comic_i_img > a > img").AttrOr("src", "http://xxx.pic"),
	}
	dinfo.Find("ul.comic_deCon_liO > li").Each(func(i int, selection *goquery.Selection) {
		data := selection.Text()
		switch i {
		case 0:
			data = strings.Replace(data, "作者：", "", 1)
			comicInfo.author = data
		case 1:
			comicInfo.isFinish = (data == "状态：已完结")
		case 2:
			data = strings.Replace(data, "类别：", "", 1)
			comicInfo.classification = data
		case 3:
			data = strings.Replace(data, "类型：", "", 1)
			comicInfo.tag = strings.Split(data, "|")
		}
	})

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterDmzj(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
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

	return comicChapter, errors.New("nil func")
}

func getChapterImageDmzj(doc *goquery.Document) (imageUrl []string, err error) {
	//fmt.Print(doc.Text())
	re := regexp.MustCompile("eval\\(.*\\)")
	picJs := re.FindString(doc.Text())
	//log.Info(picJs)
	picJson, err := evalDecode(picJs)
	picJson, err = evalDecode(picJson)
	re = regexp.MustCompile("var pages='(.*)'")
	picJsons := re.FindStringSubmatch(picJson)
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match!")
	}
	picJson = picJsons[1]
	var f map[string]interface{}
	err = json.Unmarshal([]byte(picJson), &f)
	imageUrl = strings.Split(f["page_url"].(string), "\r\n")
	return imageUrl, nil
}
