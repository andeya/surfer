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
	//phantomjs内核
	jar, _ := cookiejar.New(nil)
	var cookies []*http.Cookie
	cookie := &http.Cookie{
		Name:   "k3",
		Value:  "v3",
		Path:   "/",
		Domain: "httpbin.org",
	}
	cookies = append(cookies, cookie)

	u, _ := url.Parse("http://httpbin.org/cookies")
	jar.SetCookies(u, cookies)

	//查看cookie------------------------------------------------------------------
	log.Println("查看cookie" + HR)
	phantomSurfer := surfer.NewPhantom("E:\\Workspace\\go-labs\\src\\lab089\\phantomjs.exe", "./tmp", jar)
	resp, err := phantomSurfer.Download(&surfer.Request{
		Url:          "http://httpbin.org/cookies",
		EnableCookie: true,
		DownloaderID: 1,
	})
	handleError(err)

	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	//设置cookie------------------------------------------------------------------
	log.Println("设置cookie" + HR)
	resp, err = phantomSurfer.Download(&surfer.Request{
		Url:          "http://httpbin.org/cookies/set?k2=v2&k1=v1",
		EnableCookie: true,
		DownloaderID: 1,
	})
	handleError(err)

	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
