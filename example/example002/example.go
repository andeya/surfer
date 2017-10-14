package main

import (
	"github.com/henrylee2cn/surfer"
	"io/ioutil"
	"log"
	"net/url"
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

	// 1.surf内核GET下载测试开始
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

	// 2.surf内核POST下载测试开始
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

	// 3.phantomjs内核GET下载测试开始
	log.Println("phantomjs内核GET下载测试开始" + HR)
	surfer.SetPhantomJsFilePath("E:/Workspace/go-labs/src/lab089/lab003/phantomjs/phantomjs.exe") // 指定phantomjs可执行文件的位置，绝对路径
	resp, err = surfer.Download(&surfer.Request{
		Url:            "http://httpbin.org/get",
		DownloaderID:   1,
		PhantomTimeout: time.Millisecond * 2000, //设置超时时间
	})

	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	// 4.phantomjs内核POST下载测试开始---------------------------------------------------------------------------
	surfer.SetPhantomJsFilePath("../phantomjs.exe") //相对路径
	log.Println("phantomjs内核POST下载测试开始" + HR)
	resp, err = surfer.Download(&surfer.Request{
		DownloaderID:   1,
		Url:            "http://httpbin.org/post",
		Method:         "POST",
		Body:           form,
		PhantomTimeout: time.Millisecond * 2000, //设置超时时间
	})
	handleError(err)
	log.Println("resp.Status=", resp.Status)
	log.Println("resp.Header=", resp.Header)

	b, err = ioutil.ReadAll(resp.Body)
	handleError(err)
	resp.Body.Close()
	log.Println("resp.Body=", string(b))

	surfer.DestroyJsFiles()

	time.Sleep(10e9)
}

func handleError(err error) {
	if err != nil {
		log.Fatal("surfer download error:", err)
	}
}
