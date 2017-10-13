package main

import (
	"github.com/henrylee2cn/surfer"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"time"
)

const (
	HR = "--------------------------------------------------------------------------"
)

func main() {
	var values, _ = url.ParseQuery("username=123456@qq.com&password=123456&login_btn=login_btn&submit=login_btn")
	log.Println("values:", values)
	var form = surfer.Form{
		Values: values,
	}

	// surf内核GET下载测试开始---------------------------------------------------------------------------
	log.Println("surf内核GET下载测试开始" + HR)
	resp, err := surfer.Download(&surfer.Request{
		Url: "http://httpbin.org/get",
	})
	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err := ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	// surf内核POST下载测试开始---------------------------------------------------------------------------
	log.Println("surf内核POST下载测试开始" + HR)
	req := &surfer.Request{
		Url:    "http://httpbin.org/post",
		Method: "POST",
		Body:   form,
	}
	b, err = req.ReadBody()
	log.Println("req.Body=", string(b))

	resp, err = surfer.Download(req)
	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	// phantomjs内核GET下载测试开始---------------------------------------------------------------------------
	log.Println("phantomjs内核GET下载测试开始" + HR)
	// 指定phantomjs可执行文件的位置
	surfer.SetPhantomJsFilePath("E:/Workspace/go-labs/src/lab089/lab003/phantomjs/phantomjs.exe")

	resp, err = surfer.Download(&surfer.Request{
		Url:          "http://httpbin.org/get",
		DownloaderID: 1,
	})
	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	os.Exit(2)

	//phantomjs内核POST下载测试开始---------------------------------------------------------------------------
	log.Println("phantomjs内核POST下载测试开始" + HR)
	// 指定使用phantomjs内核下载
	resp, err = surfer.Download(&surfer.Request{
		DownloaderID: 1,
		Url:          "http://accounts.lewaos.com/",
		Method:       "POST",
		Body:         form,
	})
	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", resp.Body)

	surfer.DestroyJsFiles()

	time.Sleep(10e9)
}

func handleError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
