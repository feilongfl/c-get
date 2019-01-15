package source

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
	"strconv"
	"strings"
)

var ParseHHXXEE = Parse_s{
	name:               "99770",
	Id:                 0,
	regex:              []string{"99770.hhxxee.com"},
	getComicInfoReq:    getComicInfoReqDefault,
	getComicInfo:       getComicInfoHhxxee,
	getComicChapterReq: getComicInfoReqDefault,
	getComicChapter:    getComicChapterHhxxee,
	getChapterImageReq: getChapterImageReqDefault,
	getChapterImage:    getChapterImageHhxxee,
	getImage:           getImageDefault,
}

func getComicInfoHhxxee(doc *goquery.Document) (comicInfo comicInfo_s, err error) {
	comicInfo = comicInfo_s{
		title:    strings.Split(strings.Trim(doc.Find("h1.cTitle").Text(), " \n"), "\n")[0],
		coverUrl: doc.Find("div.cDefaultImg > img").AttrOr("src", _unknowPic_),
		info:     strings.Trim(doc.Find("div.cCon").Text(), " \n"),
		isFinish: false,
	}

	log.WithFields(log.Fields{
		"comic": comicInfo,
	}).Info("info-comic")

	return comicInfo, nil
}

func getComicChapterHhxxee(doc *goquery.Document) (comicChapter []comicChapter_s, err error) {
	comicChapter = make([]comicChapter_s, 0)
	doc.Find(".cVolList > div").Each(func(i int, selection *goquery.Selection) {
		adoc := selection.Find("a")

		var c = comicChapter_s{
			name:  adoc.Text(),
			group: 0,
		}

		c.url = adoc.AttrOr("href", "")
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

var serverList = []string{
	"http://165.94201314.net/dm01/",
	"http://165.94201314.net/dm02/",
	"http://165.94201314.net/dm03/",
	"http://165.94201314.net/dm04/",
	"http://165.94201314.net/dm05/",
	"http://165.94201314.net/dm06/",
	"http://165.94201314.net/dm07/",
	"http://165.94201314.net/dm08/",
	"http://165.94201314.net/dm09/",
	"http://165.94201314.net/dm10/",
	"http://165.94201314.net/dm11/",
	"http://165.94201314.net/dm12/",
	"http://165.94201314.net/dm13/",
	"http://173.231.57.238/dm14/",
	"http://165.94201314.net/dm15/",
	"http://142.4.34.102/dm16/",
}

func hhxxeeServerGet(url string) (serverId int) {
	re := regexp.MustCompile(`ok-comic(\d+)`)
	_serverId, err := strconv.ParseInt(re.FindStringSubmatch(url)[1], 10, 32)
	if err != nil {
		return 0
	}
	return int(_serverId - 1)
}

func getChapterImageHhxxee(doc *goquery.Document) (imageUrl []string, err error) {
	re := regexp.MustCompile(`var sFiles="(.*?)"`)
	a := doc.Text()
	picJsons := re.FindStringSubmatch(a)
	if len(picJsons) != 2 {
		return imageUrl, errors.New("regex not match")
	}
	decoded := strings.Split(picJsons[1], "|")
	for _, p := range decoded {
		imageUrl = append(imageUrl, serverList[hhxxeeServerGet(p)]+p)
	}

	return imageUrl, nil
}
