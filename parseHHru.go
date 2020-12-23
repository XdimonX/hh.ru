package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
	cookieJar http.CookieJar
	client    *http.Client
)

func authorization() {
	var err error
	cookieJar, err = cookiejar.New(nil)
	checkErr(err)
	client = &http.Client{Jar: cookieJar}
	response, err := client.Get("https://togliatti.hh.ru/account/login")
	urll, _ := url.Parse("https://togliatti.hh.ru/account/login")
	cookie := cookieJar.Cookies(urll)
	xsfr := ""
	for _, v := range cookie {
		if v.Name == "_xsrf" {
			xsfr = v.Value
			break
		}
	}

	data := url.Values{
		"_xsrf":    {xsfr},
		"backUrl":  {"https://togliatti.hh.ru/"},
		"failUrl":  {"/account/login?backurl=%%2F"},
		"remember": {"false"},
		"username": {loginHHru},
		"password": {passwordHHru},
		"isBot":    {"false"},
	}
	buf := &bytes.Buffer{}
	writer := multipart.NewWriter(buf)
	writer.CreatePart()
	req, _ := http.NewRequest("POST", "https://togliatti.hh.ru/account/login", strings.NewReader(data.Encode()))
	req.Header.Set("accept", "application/json")
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 OPR/72.0.3815.454")
	req.Header.Set("sec-fetch-dest", "empty")
	req.Header.Set("sec-fetch-mode", "cors")
	req.Header.Set("sec-fetch-site", "same-origin")
	req.Header.Set("content-type", "multipart/form-data; boundary=----WebKitFormBoundaryrC1sHenhCYf93NTr")
	req.Header.Set("content-length", "820")
	req.Header.Set("origin", "https://togliatti.hh.ru")
	req.Header.Set("referer", "https://togliatti.hh.ru/account/login")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("x-hhtmsource", "account_login")
	req.Header.Set("x-requested-with", "XMLHttpRequest")
	req.Header.Set("x-xsrftoken", xsfr)
	response, err = client.Do(req)
	// response, err = client.PostForm("https://togliatti.hh.ru/account/login?backurl=%%2F", url.Values{
	// 	"_xsrf":    {xsfr},
	// 	"backUrl":  {"https://togliatti.hh.ru/"},
	// 	"failUrl":  {"/account/login?backurl=%%2F"},
	// 	"remember": {"false"},
	// 	"username": {loginHHru},
	// 	"password": {passwordHHru},
	// 	"isBot":    {"false"},
	// })
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
	ioutil.WriteFile("out.html", body, os.ModeExclusive)
	checkErr(err)
	if response.StatusCode != 200 {
		os.Exit(2)
	}

	response.Body.Close()
	fmt.Println(string(body))

}

func getResumeList() {
	authorization()
	req, _ := http.NewRequest("GET", "https://togliatti.hh.ru/applicant/resumes", nil)
	req.Header.Set("accept-encoding", "gzip, deflate, br")
	req.Header.Set("accept-language", "ru-RU,ru;q=0.9,en-US;q=0.8,en;q=0.7")
	req.Header.Set("user-agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/86.0.4240.198 Safari/537.36 OPR/72.0.3815.454")
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	// req.Proto = "HTTP/2.0"
	// req.ProtoMajor = 2
	// req.ProtoMinor = 0
	// response, err := client.Get("https://togliatti.hh.ru/applicant/resumes")
	response, err := client.Do(req)
	checkErr(err)
	cookie := response.Cookies()
	_ = cookie
	urll, _ := url.Parse("https://togliatti.hh.ru/applicant/resumes")
	cookie = client.Jar.Cookies(urll)

	doc, err := goquery.NewDocumentFromReader(response.Body)
	var p []byte
	response.Body.Read(p)
	println(p)
	checkErr(err)
	doc.Find(".bloko-column_l-11").Each(func(i int, s *goquery.Selection) {
		println(i)
	})

}
