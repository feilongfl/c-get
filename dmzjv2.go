package main

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var ParseDMZJV2 = parse_s{
	name:               "动漫之家v2",
	id:                 0,
	regex:              []string{"manhua.dmzj.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoDmzjV2,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterDmzjV2,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImageDmzjV2,
	getImage:           getImageDefault,
}

func getComicInfoDmzjV2(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	//dinfo := doc.Find("div.wrap_intro_l")
	comicInfo = comicInfo_s{
		title:          doc.Find("span.anim_title_text > a > h1").Text(),
		coverUrl:       doc.Find("div.comic_i_img > a > img").AttrOr("src", _unknowPic_),
		info:           doc.Find("div.line_height_content").Text(),
		author:         doc.Find("tr:nth-child(3) > td > a").Text(),
		classification: doc.Find("tr:nth-child(7) > td > a").Text(),
		isFinish:       doc.Find("tr:nth-child(5) > td > a").Text() != "连载中",
		//tag:            doc.Find("div.line_height_content").Text(),
	}

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterDmzjV2(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	comicChapter = make([]comicChapter_s, 0)
	//docS := doc.Find("div.zj_list_con").First()
	doc.Find("div.cartoon_online_border > ul > li").Each(func(i int, selection *goquery.Selection) {
		adoc := selection.Find("a")
		var c = comicChapter_s{
			name:  adoc.AttrOr("title", ""),
			url:   "https://manhua.dmzj.com" + adoc.AttrOr("href", ""),
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

func getChapterImageDmzjV2(doc *goquery.Document) (imageUrl []string, err error) {
	//fmt.Print(doc.Text())
	re := regexp.MustCompile("eval\\(.*\\)")
	picJs := re.FindString(doc.Text())
	log.Info(picJs)
	picJson, err := evalDecode(picJs)
	//log.Info(picJson)
	//picJson, err = evalDecode(picJson)
	re = regexp.MustCompile(`var pages=pages='\[(.*)\]'`)
	picJsons := re.FindStringSubmatch(picJson)
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match")
	}
	picJson = picJsons[1]
	//var f map[string]interface{}
	//err = json.Unmarshal([]byte(picJson), &f)
	picJson = strings.Replace(picJson, `"`, ``, -1)
	picJson = strings.Replace(picJson, "\\", ``, -1)
	imageUrl = strings.Split(picJson, ",")
	for index, url := range imageUrl {
		imageUrl[index] = "https://images.dmzj.com/" + url
	}
	return imageUrl, nil
}
