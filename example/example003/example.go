package main

import (
	"github.com/henrylee2cn/surfer"
	"log"
	"time"
)

const (
	HR = "--------------------------------------------------------------------------"
)

func main() {
	//默认内核
	log.Println("默认内核timeout" + HR)
	log.Println(time.Now(), "start")
	resp, err := surfer.Download(&surfer.Request{
		Url:         "https://www.google.com/",
		DialTimeout: time.Second * 1,
		TryTimes:    2,
	})
	log.Println(time.Now(), "timeout")
	if err != nil {
		log.Println("surfer download error:", err)
	} else {
		log.Println("resp.Status=", resp.Status)
	}

	//phantomjs内核
	surfer.SetPhantomJsFilePath("E:/Workspace/go-labs/src/lab089/lab003/phantomjs/phantomjs.exe")
	log.Println("phantomjs内核" + HR)
	log.Println(time.Now(), "start")
	resp, err = surfer.Download(&surfer.Request{
		Url:          "https://www.google.com",
		DialTimeout:  time.Second * 2,
		DownloaderID: 1,
		TryTimes:     3,
	})
	log.Println(time.Now(), "timeout")
	if err != nil {
		log.Println("surfer download error:", err)
	} else {
		log.Println("resp.Status=", resp.Status)
	}

	surfer.DestroyJsFiles()
	time.Sleep(10e9)
}
