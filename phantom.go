// Copyright 2015 henrylee2cn Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package surfer

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type (
	// Phantom 基于Phantomjs的下载器实现，作为surfer的补充
	// 效率较surfer会慢很多，但是因为模拟浏览器，破防性更好
	// 支持UserAgent/TryTimes/RetryPause/自定义js
	Phantom struct {
		PhantomjsFile string            //Phantomjs完整文件名
		TempJsDir     string            //临时js存放目录
		jsFileMap     map[string]string //已存在的js文件
		CookieJar     *cookiejar.Jar
	}
	// Response 用于解析Phantomjs的响应内容
	Response struct {
		Cookies []string
		Body    string
	}
)

// NewPhantom 创建一个Phantomjs下载器
func NewPhantom(phantomjsFile, tempJsDir string, jar ...*cookiejar.Jar) Surfer {
	phantom := &Phantom{
		PhantomjsFile: phantomjsFile,
		TempJsDir:     tempJsDir,
		jsFileMap:     make(map[string]string),
	}
	if len(jar) != 0 {
		phantom.CookieJar = jar[0]
	} else {
		phantom.CookieJar, _ = cookiejar.New(nil)
	}
	if !filepath.IsAbs(phantom.PhantomjsFile) {
		phantom.PhantomjsFile, _ = filepath.Abs(phantom.PhantomjsFile)
	}
	if !filepath.IsAbs(phantom.TempJsDir) {
		phantom.TempJsDir, _ = filepath.Abs(phantom.TempJsDir)
	}
	// 创建/打开目录
	err := os.MkdirAll(phantom.TempJsDir, 0777)
	if err != nil {
		log.Printf("[E] Surfer: %v\n", err)
		return phantom
	}
	phantom.createJsFile("js", js)
	return phantom
}

// Download 实现surfer下载器接口
func (phantom *Phantom) Download(req *Request) (resp *http.Response, err error) {
	err = req.prepare()
	if err != nil {
		return resp, err
	}
	var encoding = "utf-8"
	if _, params, err := mime.ParseMediaType(req.Header.Get("Content-Type")); err == nil {
		if cs, ok := params["charset"]; ok {
			encoding = strings.ToLower(strings.TrimSpace(cs))
		}
	}

	req.Header.Del("Content-Type")

	if req.EnableCookie {
		log.Println("req.EnableCookie:", req.EnableCookie)
		_req := http.Request{Header: req.Header}
		for _, cookie := range phantom.CookieJar.Cookies(req.url) {
			log.Println("liguoqinjim 发起请求前添加cookie:", cookie)
			_req.AddCookie(cookie)
		}
	}

	var b, _ = req.ReadBody()

	//todo writeback里面的重置url会引起Bug
	resp = req.writeback(resp)

	var args = []string{
		phantom.jsFileMap["js"],
		req.Url,
		req.Header.Get("Cookie"),
		encoding,
		req.Header.Get("User-Agent"),
		string(b),
		strings.ToLower(req.Method),
		fmt.Sprint(int(req.DialTimeout / time.Millisecond)),
	}

	for i := 0; i < req.TryTimes; i++ {
		if i != 0 {
			time.Sleep(req.RetryPause)
		}

		cmd := exec.Command(phantom.PhantomjsFile, args...)
		if resp.Body, err = cmd.StdoutPipe(); err != nil {
			continue
		}
		err = cmd.Start()
		if err != nil || resp.Body == nil {
			continue
		}
		var b []byte
		b, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Println("surfer read phantomjs out error:", err)
			continue
		}
		retResp := Response{}
		err = json.Unmarshal(b, &retResp)
		if err != nil {
			log.Println("surfer phantomjs out:", string(b))
			continue
		}
		resp.Header = req.Header
		resp.Header.Del("Set-Cookie")
		for _, c := range retResp.Cookies {
			resp.Header.Add("Set-Cookie", c)
		}
		if req.EnableCookie {
			if rc := resp.Cookies(); len(rc) > 0 {
				phantom.CookieJar.SetCookies(req.url, rc)
			}
		}
		resp.Body = ioutil.NopCloser(strings.NewReader(retResp.Body))
		break
	}

	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
	} else {
		resp.StatusCode = http.StatusBadGateway
		resp.Status = err.Error()
	}

	log.Printf("liguoqinjim 访问后cookies:%+v", phantom.CookieJar.Cookies(req.url))

	return resp, err
}

// DestroyJsFiles 销毁js临时文件
func (phantom *Phantom) DestroyJsFiles() {
	p, _ := filepath.Split(phantom.TempJsDir)
	if p == "" {
		return
	}
	for _, filename := range phantom.jsFileMap {
		os.Remove(filename)
	}
	if len(WalkDir(p)) == 1 {
		os.Remove(p)
	}
}

func (phantom *Phantom) createJsFile(fileName, jsCode string) {
	fullFileName := filepath.Join(phantom.TempJsDir, fileName)
	// 创建并写入文件
	f, _ := os.Create(fullFileName)
	f.Write([]byte(jsCode))
	f.Close()
	phantom.jsFileMap[fileName] = fullFileName
}

/*
* system.args[0] == js
* system.args[1] == url
* system.args[2] == cookie
* system.args[3] == pageEncode
* system.args[4] == userAgent
* system.args[5] == postdata
* system.args[6] == method
* system.args[7] == timeout
 */
const js string = `
var system = require('system');
var page = require('webpage').create();
var url = system.args[1];
var cookie = system.args[2];
var pageEncode = system.args[3];
var userAgent = system.args[4];
var postdata = system.args[5];
var method = system.args[6];
var timeout = system.args[7];
var ret = "";
var exit = function () {
  console.log(ret);
  phantom.exit();
};

phantom.outputEncoding = pageEncode;
page.settings.userAgent = userAgent;
page.settings.resourceTimeout = timeout;
page.settings.XSSAuditingEnabled = true;
page.onResourceRequested = function (requestData, request) {
  request.setHeader('Cookie', cookie)
};
page.onError = function (msg, trace) {
  console.log("error:" + msg);
};
page.onResourceTimeout = function (e) {
  console.log("phantomjs onResourceTimeout error");
  // console.log(e.errorCode);   // it'll probably be 408
  // console.log(e.errorString); // it'll probably be 'Network timeout on resource'
  // console.log(e.url);         // the url whose request timed out
  phantom.exit(1);
};
page.onResourceError = function (resourceError) {
};
page.onLoadFinished = function (status) {
  if (status !== 'success') {
    console.log("phantomjs status:" + status);
    exit();
  } else {
    var cookies = new Array();
    for (var i in page.cookies) {
      var cookie = page.cookies[i];
      var c = cookie["name"] + "=" + cookie["value"];
      for (var obj in cookie) {
        if (obj == 'name' || obj == 'value') {
          continue;
        }
        c += "; " + obj + "=" + cookie[obj];
      }
      cookies[i] = c;
    }

    var resp = {
      "Cookies": cookies,
      "Body": page.content
    };

    if (page.content.indexOf("body") != -1) {
      ret = JSON.stringify(resp);
      exit();
    }
  }
};

page.open(url, method, postdata, function (status) {
});
`
