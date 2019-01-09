package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/cheggaaa/pb.v2"
	"os"
	"runtime"
	"time"
)

type comicInfo_s struct {
	title           string
	author          string
	isFinish        bool
	lastUpdateData  string
	classification  string
	tag             []string
	coverUrl        string
	comicChapterUrl string
	info            string
}

type comicChapter_s struct {
	url     string
	name    string
	group   int
	picsUrl []string
}

type comic_s struct {
	comicInfoUrl string
	comicInfo    comicInfo_s
	comicChapter []comicChapter_s
	parse        parse_s
}

const sthreadMax = 2
const ssleepTime = 500
const threadMax = 2
const sleepTime = 1000

func parseGet(url string) (p *parse_s, err error) {
	parse, err := newParseFromUrl(url)
	if err != nil {
		log.WithFields(log.Fields{
			"url": url,
		}).Fatal(err)
	}
	log.WithFields(log.Fields{
		"parse": parse.name,
		"id":    parse.id,
	}).Info("parse match")
	return parse, err
}

func infoGet(p *parse_s, comic *comic_s) (err error) {
	doc, err := p.getComicInfoReq(comic.comicInfoUrl)
	if err != nil {
		return err
	}
	//fmt.Print(doc)
	comic.comicInfo, err = p.getComicInfo(doc)
	if comic.comicInfo.comicChapterUrl != "" {
		if doc, err = p.getComicInfoReq(comic.comicInfoUrl); err != nil {
			return err
		}
	} else {
		comic.comicInfo.comicChapterUrl = comic.comicInfoUrl
	}
	log.WithFields(log.Fields{
		"title": comic.comicInfo.title,
	}).Warning("comic get")

	comic.comicChapter, err = p.getComicChapter(doc)

	return err
}

/////////////////////////////////////////
//get all images url
func getImageUrlList(p *parse_s, comic *comic_s) (err error) {
	type picChan_s struct {
		index int
		pics  []string
	}
	c := make(chan picChan_s)
	bar := pb.StartNew(len(comic.comicChapter))
	sworker := 0
	sindex := 0

	for {
		if sworker < sthreadMax && sindex < len(comic.comicChapter) {
			time.Sleep(ssleepTime * time.Microsecond)
			sworker++
			chapter := comic.comicChapter[sindex]
			go func(chapter comicChapter_s, index int) {
				log.WithFields(log.Fields{
					"index":   index,
					"chapter": chapter.name,
				}).Info("downloading")

				doc, err := p.getChapterImageReq(chapter.url, comic.comicInfo.comicChapterUrl)
				// todo: fix here
				pics, err := p.getChapterImage(doc)
				if err != nil {
					log.WithFields(log.Fields{
						"index":       index,
						"chapter":     chapter.name,
						"chapter url": chapter.url,
					}).Warning(err)
				}
				c <- picChan_s{index, pics}

				sindex++
				runtime.Gosched()
			}(chapter, sindex)
		} else {
			done := <-c
			comic.comicChapter[done.index].picsUrl = done.pics
			log.WithFields(log.Fields{
				"done": done,
			}).Info("done")
			bar.Increment()
			sworker--
			if bar.Current() == int64(len(comic.comicChapter)) {
				break
			}
		}
	}
	bar.Finish()

	log.WithFields(log.Fields{
		"comic.comicChapter": comic.comicChapter,
	}).Debug("done all")

	return err
}

/////////////////////////////////////////
//get all images
func downloadImages(p *parse_s, comic *comic_s) (err error) {
	type imageDownload_s struct {
		index       int
		chapterId   int
		chapterName string
		chapterUrl  string
		imageUrl    string
		savepath    string
		success     bool
		retry       int
	}
	imageDown_c := make(chan imageDownload_s)
	//threadWork := make(chan int)

	imageDownloadList := make([]imageDownload_s, 0)
	imageid := 0
	for chapterindex, chapter := range comic.comicChapter {
		os.MkdirAll(fmt.Sprintf("./data/%v/%v", comic.comicInfo.title, chapter.name), os.ModePerm)
		for imageindex, image := range chapter.picsUrl {
			imageDownloadList = append(imageDownloadList, imageDownload_s{
				imageid,
				chapterindex,
				chapter.name,
				chapter.url,
				image,
				fmt.Sprintf("./data/%v/%v/%v.jpg", comic.comicInfo.title, chapter.name, imageindex),
				false,
				0,
			})
			imageid += 1
		}
	}
	log.WithFields(log.Fields{
		"imageDownloadList": imageDownloadList,
	}).Debug("image all")

	bar := pb.StartNew(len(imageDownloadList))
	download := func(image imageDownload_s) {
		log.WithFields(log.Fields{
			"image": image,
		}).Info("start:")

		if err := p.getImage(image.imageUrl, image.chapterUrl, image.savepath); err == nil {
			image.success = true
		}
		imageDown_c <- image
		runtime.Gosched()
	}
	downloadDone := func(done imageDownload_s) {
		time.Sleep(sleepTime * time.Microsecond)

		if done.success != true {
			done.retry += 1
			if done.retry < 5 {
				log.WithFields(log.Fields{
					"image": done,
				}).Warning("retry:")
				go download(done)
			} else {
				log.WithFields(log.Fields{
					"image": done,
				}).Error("failed:")
				bar.Increment()
			}
		} else {
			bar.Increment()
			//if bar.Current() == int64(len(imageDownloadList)) {
			//	break
			//}
		}
	}

	works := 0
	//for _, image := range imageDownloadList {
	for {
		if works < threadMax && bar.Current() < int64(len(imageDownloadList)) {
			works++
			image := imageDownloadList[bar.Current()]
			go download(image)
		} else {
			done := <-imageDown_c
			works--
			downloadDone(done)
			if bar.Current() == int64(len(imageDownloadList)) {
				break
			}
		}
	}
	//for {
	//	done := <-imageDown_c
	//	downloadDone(done)
	//}
	bar.Finish()

	return err
}

func infoComic(url string) (err error) {
	comic := comic_s{comicInfoUrl: url}
	parse, err := parseGet(url)
	if err != nil {
		return err
	}

	if err = infoGet(parse, &comic); err != nil {
		return err
	}

	if err = getImageUrlList(parse, &comic); err != nil {
		return err
	}

	return downloadImages(parse, &comic)
}
