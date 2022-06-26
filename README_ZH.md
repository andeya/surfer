# Surfer [![GitHub release](https://img.shields.io/github/release/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/releases) [![report card](https://goreportcard.com/badge/github.com/andeya/surfer?style=flat-square)](http://goreportcard.com/report/andeya/surfer) [![github issues](https://img.shields.io/github/issues/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/issues?q=is%3Aopen+is%3Aissue) [![github closed issues](https://img.shields.io/github/issues-closed-raw/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/issues?q=is%3Aissue+is%3Aclosed) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/andeya/surfer) [![view Go大数据](https://img.shields.io/badge/官方QQ群-Go大数据(42731170)-27a5ea.svg?style=flat-square)](http://jq.qq.com/?_wv=1027&k=XnGGnc)


Surfer 是一款Go语言编写的高并发 web 客户端，拥有surf与phantom两种下载内核，高度模拟浏览器行为，可实现模拟登录等功能。

高并发爬虫[Pholcus](https://github.com/andeya/pholcus)的专用下载器。

## 特性

- 支持 `surf` 和 `phantomjs` 两种下载内核
- 支持大量随机的User-Agent
- 支持缓存cookie
- 支持`http`/`https`两种协议

## 用法
```
package main

import (
    "github.com/andeya/surfer"
    "io/ioutil"
    "log"
)

func main() {
    // 默认使用surf内核下载
    resp, err := surfer.Download(&surfer.Request{
        Url: "http://github.com/andeya/surfer",
    })
    if err != nil {
        log.Fatal(err)
    }
    b, err := ioutil.ReadAll(resp.Body)
    log.Println(string(b), err)

    // 指定使用phantomjs内核下载
    surfer.SetPhantomJsFilePath("Path to phantomjs.exe")
    resp, err = surfer.Download(&surfer.Request{
        Url:          "http://github.com/andeya",
        DownloaderID: 1,
    })
    if err != nil {
        log.Fatal(err)
    }
    b, err = ioutil.ReadAll(resp.Body)
    log.Println(string(b), err)

    resp.Body.Close()
    surfer.DestroyJsFiles()
}
```

[完整示例](https://github.com/andeya/surfer/tree/master/example)


## 开源协议

Surfer 项目采用商业应用友好的[Apache License v2](https://github.com/andeya/surfer/raw/master/LICENSE).发布
