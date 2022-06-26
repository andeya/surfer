# Surfer [![GitHub release](https://img.shields.io/github/release/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/releases) [![report card](https://goreportcard.com/badge/github.com/andeya/surfer?style=flat-square)](http://goreportcard.com/report/andeya/surfer) [![github issues](https://img.shields.io/github/issues/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/issues?q=is%3Aopen+is%3Aissue) [![github closed issues](https://img.shields.io/github/issues-closed-raw/andeya/surfer.svg?style=flat-square)](https://github.com/andeya/surfer/issues?q=is%3Aissue+is%3Aclosed) [![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/andeya/surfer) [![view Go大数据](https://img.shields.io/badge/官方QQ群-Go大数据(42731170)-27a5ea.svg?style=flat-square)](http://jq.qq.com/?_wv=1027&k=XnGGnc)

Package surfer is a high level concurrency http client.
It has `surf` and` phantom` download engines, highly simulated browser behavior, the function of analog login and so on.

[简体中文](https://github.com/andeya/surfer/blob/master/README_ZH.md)

## Features
- Both `surf` and `phantomjs` engines are supported
- Support random User-Agent
- Support cache cookie
- Support http/https

## Usage
```
package main

import (
    "github.com/andeya/surfer"
    "io/ioutil"
    "log"
)

func main() {
    // Use surf engine
    resp, err := surfer.Download(&surfer.Request{
        Url: "http://github.com/andeya/surfer",
    })
    if err != nil {
        log.Fatal(err)
    }
    b, err := ioutil.ReadAll(resp.Body)
    log.Println(string(b), err)

    // Use phantomjs engine
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
[Full example](https://github.com/andeya/surfer/tree/master/example)

## License

Surfer is under Apache v2 License. See the [LICENSE](https://github.com/andeya/surfer/raw/master/LICENSE) file for the full license text.
