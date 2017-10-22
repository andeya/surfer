package main

import (
	"github.com/henrylee2cn/surfer"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
)

const (
	HR = "------------------------------------------------------------------"
)

func main() {
	//默认内核
	jar1, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "k3",
		Value:  "v3",
		Path:   "/",
		Domain: "httpbin.org",
	}
	cookies = append(cookies, cookie)

	u, _ := url.Parse("http://httpbin.org/cookies")
	jar1.SetCookies(u, cookies)

	//查看cookie
	log.Println("查看Cookie" + HR)
	defaultSurfer := surfer.New(jar1)
	resp, err := defaultSurfer.Download(&surfer.Request{
		Url:          "http://httpbin.org/cookies",
		EnableCookie: true,
	})
	handleError(err)

	log.Println("resp.Status=", resp.Status)
	b, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println("body=", string(b))
	log.Println("cookie=", jar1.Cookies(u))

	//设置cookie
	log.Println("设置Cookie" + HR)
	resp, err = defaultSurfer.Download(&surfer.Request{
		Url:          "http://httpbin.org/cookies/set?k2=v2&k1=v1",
		EnableCookie: true,
	})
	handleError(err)

	log.Println("resp.Status=", resp.Status)
	b, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	log.Println("body=", string(b))

	log.Println(jar1.Cookies(u))
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
