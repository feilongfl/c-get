package source

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"regexp"
)

type Parse_s struct {
	name               string
	Id                 int
	regex              []string
	getComicInfoReq    func(url string) (*goquery.Document, error)
	getComicInfo       func(doc *goquery.Document) (comicInfo comicInfo_s, err error)
	getComicChapterReq func(url string) (*goquery.Document, error)
	getComicChapter    func(doc *goquery.Document) (comicChapter []comicChapter_s, err error)
	getChapterImageReq func(url string, referer string) (*goquery.Document, error)
	getChapterImage    func(doc *goquery.Document) (imageUrl []string, err error)
	getImage           func(imageUrl string, referer string, path string) (err error)
}

var ParseList = make([]Parse_s, 0)

func newParseFromUrl(url string) (p *Parse_s, err error) {
	for _, p := range ParseList {
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
