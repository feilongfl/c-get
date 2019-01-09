package core

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"net/url"
	"os"

	"net/http"
)

type HttpReqStruct struct {
	Url       string
	Referer   string
	UserAgent string
	Method    string
	body      io.Reader
}

var defaultHttpReq = HttpReqStruct{
	Method:    "GET",
	UserAgent: "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.110 Safari/537.36",
	body:      nil,
	Referer:   "",
}

func DefaultHttpReq(url string) (httpReq HttpReqStruct) {
	httpReq = defaultHttpReq
	httpReq.Url = url
	return httpReq
}

func ProxyAwareHttpClient() *http.Client {
	httpTransport := &http.Transport{}
	proxyServer, isSet := os.LookupEnv("HTTP_PROXY")

	if !isSet {
		log.Info("no proxy")
		httpClient := &http.Client{Transport: httpTransport}
		return httpClient
	}

	proxyUrl, err := url.Parse(proxyServer)
	if err != nil {
		log.WithFields(
			log.Fields{"proxyUrl": proxyUrl},
		).Warning("proxy is invalid")
	} else {
		switch proxyUrl.Scheme {
		case "http":
			httpTransport.Proxy = http.ProxyURL(proxyUrl)
			log.WithFields(
				log.Fields{"proxyUrl": proxyUrl},
			).Info(proxyUrl.Scheme + " proxy is set")
		default:
			log.WithFields(
				log.Fields{"proxyUrl": proxyUrl},
			).Warning(proxyUrl.Scheme + " proxy not support")
			break
		}
	}

	httpClient := &http.Client{Transport: httpTransport}
	return httpClient
}

func httpGet(httpReq HttpReqStruct) ([]byte, error) {
	log.WithFields(log.Fields{
		"req": httpReq,
	}).Debug("httpGet")

	client := ProxyAwareHttpClient()

	req, err := http.NewRequest(httpReq.Method, httpReq.Url, httpReq.body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	req.Header.Set("User-Agent", httpReq.UserAgent)
	if httpReq.Referer != "" {
		req.Header.Set("Referer", httpReq.Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	//log.Println(string(body))
	return body, nil
}

func HtmlGet(httpReq HttpReqStruct) (*goquery.Document, error) {
	log.WithFields(log.Fields{
		"req": httpReq,
	}).Debug("httpGet")

	client := ProxyAwareHttpClient()

	req, err := http.NewRequest(httpReq.Method, httpReq.Url, httpReq.body)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	if httpReq.UserAgent != "" {
		req.Header.Set("User-Agent", httpReq.UserAgent)
	}
	if httpReq.Referer != "" {
		req.Header.Set("Referer", httpReq.Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	//log.Println(string(body))
	return doc, nil
}

func FileGet(httpReq HttpReqStruct, storagePath string) error {
	log.WithFields(log.Fields{
		"req": httpReq,
	}).Debug("FileGet")

	client := ProxyAwareHttpClient()

	req, err := http.NewRequest(httpReq.Method, httpReq.Url, httpReq.body)
	if err != nil {
		log.Error(err)
		return err
	}

	req.Header.Set("User-Agent", httpReq.UserAgent)
	if httpReq.Referer != "" {
		req.Header.Set("Referer", httpReq.Referer)
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Error(err)
		return err
	}
	if resp.StatusCode != 200 {
		return errors.New(resp.Status)
	}
	defer resp.Body.Close()

	f, err := os.Create(storagePath)
	if err != nil {
		log.Error(err)
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, resp.Body)
	if err != nil {
		log.Error(err)
		return err
	}

	return nil
}
