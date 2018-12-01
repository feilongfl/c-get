package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gopkg.in/cheggaaa/pb.v2"
	"os"
	"runtime"
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

func infoComic(url string) (err error) {
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

	var comic = comic_s{comicInfoUrl: url}
	doc, err := parse.getComicInfoReq(comic.comicInfoUrl)
	if err != nil {
		return err
	}
	//fmt.Print(doc)
	comic.comicInfo, err = parse.getComicInfo(doc)
	if comic.comicInfo.comicChapterUrl != "" {
		doc, err = parse.getComicInfoReq(comic.comicInfoUrl)
		if err != nil {
			return err
		}
	} else {
		comic.comicInfo.comicChapterUrl = comic.comicInfoUrl
	}
	log.WithFields(log.Fields{
		"title": comic.comicInfo.title,
	}).Warning("comic get")
	comic.comicChapter, err = parse.getComicChapter(doc)
	if err != nil {
		return err
	}

	//get all images url
	type picChan_s struct {
		index int
		pics  []string
	}
	c := make(chan picChan_s)
	bar := pb.StartNew(len(comic.comicChapter))
	for index, chapter := range comic.comicChapter {
		go func(chapter comicChapter_s, index int) {
			//chapter := comic.comicChapter[index]
			log.WithFields(log.Fields{
				"index":   index,
				"chapter": chapter.name,
			}).Info("downloading")

			doc, err = parse.getChapterImageReq(chapter.url, comic.comicInfo.comicChapterUrl)
			// todo: fix here
			pics, err := parse.getChapterImage(doc)
			if err != nil {
				log.WithFields(log.Fields{
					"index":       index,
					"chapter":     chapter.name,
					"chapter url": chapter.url,
				}).Warning(err)
			}
			c <- picChan_s{index, pics}
			runtime.Gosched()
		}(chapter, index)
	}
	for {
		done := <-c
		comic.comicChapter[done.index].picsUrl = done.pics
		bar.Increment()
		if bar.Current() == int64(len(comic.comicChapter)) {
			break
		}
	}
	bar.Finish()

	//get all images
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
	bar = pb.StartNew(len(imageDownloadList))
	download := func(image imageDownload_s) {
		err := parse.getImage(image.imageUrl, image.chapterUrl, image.savepath)
		if err == nil {
			image.success = true
		}
		imageDown_c <- image
		runtime.Gosched()
	}
	for _, image := range imageDownloadList {
		go download(image)
	}
	for {
		done := <-imageDown_c
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
			}
		} else {
			bar.Increment()
			if bar.Current() == int64(len(imageDownloadList)) {
				break
			}
		}
	}
	bar.Finish()

	return nil
}
