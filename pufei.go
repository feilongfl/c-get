package main

import (
	"encoding/base64"
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strings"
)

var ParsePUFEI = parse_s{
	name:               "扑飞漫画",
	id:                 0,
	regex:              []string{"www.pufei.net"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoPufei,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterPufei,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImagePufei,
	getImage:           getImageDefault,
}

func getComicInfoPufei(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	dinfo := doc.Find("div.detail")
	comicInfo = comicInfo_s{
		title:    ConvertToString(dinfo.Find("div.titleInfo > h1").Text(), "gbk", "utf-8"),
		coverUrl: dinfo.Find("div.info_cover > img").AttrOr("src", _unknowPic_),
		isFinish: !(dinfo.Find("div.titleInfo > span").Text() == "连载"),
	}
	dinfo.Find("div.detailInfo > ul > li").Each(func(i int, selection *goquery.Selection) {
		data := ConvertToString(selection.Text(), "gbk", "utf-8")
		switch i {
		case 0:
			data = strings.Replace(data, "更新时间：", "", 1)
			comicInfo.lastUpdateData = data
		case 1:
			data = strings.Replace(data, "作者：", "", 1)
			comicInfo.author = data
		case 2:
			data = strings.Replace(data, "类别：", "", 1)
			comicInfo.classification = data
		case 6:
			data = strings.Replace(data, "关键词：", "", 1)
			comicInfo.tag = strings.Split(data, ",")
		}
	})

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterPufei(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	comicChapter = make([]comicChapter_s, 0)
	//docS := doc.Find("div.plistBox > div.plist > ul > li").First()
	doc.Find("div.plistBox > div.plist > ul > li").Each(func(i int, selection *goquery.Selection) {
		adoc := selection.Find("a")

		var c = comicChapter_s{
			name: ConvertToString(adoc.AttrOr("title", ""), "gbk", "utf-8"),

			group: 0,
		}

		adochref := adoc.AttrOr("href", "")
		if adochref != "" {
			if strings.Index(adochref, "https://") == -1 {
				c.url = "http://www.pufei.net" + adochref
			} else {
				c.url = adochref
			}
			comicChapter = append(comicChapter, c)
		}
	})

	log.WithFields(log.Fields{
		"comicChapter length": len(comicChapter),
		"comicChapter":        comicChapter,
	}).Debug("info-comic")

	return comicChapter, nil
}

func getChapterImagePufei(doc *goquery.Document) (imageUrl []string, err error) {
	re := regexp.MustCompile(`packed="(.*?)"`)
	picJsons := re.FindStringSubmatch(doc.Text())
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match")
	}
	decoded, err := base64.StdEncoding.DecodeString(picJsons[1])
	if err != nil {
		return imageUrl, err
	}
	picJson, err := evalDecode(string(decoded))
	if err != nil {
		return imageUrl, err
	}
	re = regexp.MustCompile(`"(.*?)"`)
	for _, p := range re.FindAllStringSubmatch(picJson, -1) {
		imageUrl = append(imageUrl, "http://res.img.pufei.net/"+p[1])
	}

	return imageUrl, nil
}
