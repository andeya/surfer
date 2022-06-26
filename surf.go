// Copyright 2015 andeya Author. All Rights Reserved.
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
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"crypto/tls"
	"io"
	"math/rand"
	"net"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"

	"github.com/andeya/goutil"
)

// Surf is the default Download implementation.
type Surf struct {
	CookieJar *cookiejar.Jar
}

// New 创建一个Surf下载器
func New(jar ...*cookiejar.Jar) Surfer {
	s := new(Surf)
	if len(jar) != 0 {
		s.CookieJar = jar[0]
	} else {
		s.CookieJar, _ = cookiejar.New(nil)
	}
	return s
}

// Download 实现surfer下载器接口
func (surf *Surf) Download(param *Request) (*http.Response, error) {
	err := param.prepare()
	if err != nil {
		return nil, err
	}
	param.Header.Set("Connection", "close")
	param.client = surf.buildClient(param)
	resp, err := surf.httpRequest(param)

	if err == nil {
		switch resp.Header.Get("Content-Encoding") {
		case "gzip":
			var gzipReader *gzip.Reader
			gzipReader, err = gzip.NewReader(resp.Body)
			if err == nil {
				resp.Body = gzipReader
			}

		case "deflate":
			resp.Body = flate.NewReader(resp.Body)

		case "zlib":
			var readCloser io.ReadCloser
			readCloser, err = zlib.NewReader(resp.Body)
			if err == nil {
				resp.Body = readCloser
			}
		}
	}

	return param.writeback(resp), err
}

var dnsCache = &DnsCache{ipPortLib: goutil.AtomicMap()}

// DnsCache DNS cache
type DnsCache struct {
	ipPortLib goutil.Map
}

// Reg registers ipPort to DNS cache.
func (d *DnsCache) Reg(addr, ipPort string) {
	d.ipPortLib.Store(addr, ipPort)
}

// Del deletes ipPort from DNS cache.
func (d *DnsCache) Del(addr string) {
	d.ipPortLib.Delete(addr)
}

// Query queries ipPort from DNS cache.
func (d *DnsCache) Query(addr string) (string, bool) {
	ipPort, ok := d.ipPortLib.Load(addr)
	if !ok {
		return "", false
	}
	return ipPort.(string), true
}

// buildClient creates, configures, and returns a *http.Client type.
func (surf *Surf) buildClient(req *Request) *http.Client {
	client := &http.Client{
		CheckRedirect: req.checkRedirect,
	}

	if req.EnableCookie {
		client.Jar = surf.CookieJar
	}

	transport := &http.Transport{
		Dial: func(network, addr string) (net.Conn, error) {
			var (
				c          net.Conn
				err        error
				ipPort, ok = dnsCache.Query(addr)
			)
			if !ok {
				ipPort = addr
				defer func() {
					if err == nil {
						dnsCache.Reg(addr, c.RemoteAddr().String())
					}
				}()
			} else {
				defer func() {
					if err != nil {
						dnsCache.Del(addr)
					}
				}()
			}
			c, err = net.DialTimeout(network, ipPort, req.DialTimeout)
			if err != nil {
				return nil, err
			}
			if req.ConnTimeout > 0 {
				c.SetDeadline(time.Now().Add(req.ConnTimeout))
			}
			return c, nil
		},
	}

	if req.proxy != nil {
		transport.Proxy = http.ProxyURL(req.proxy)
	}

	if strings.ToLower(req.url.Scheme) == "https" {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	return client
}

// send uses the given *http.Request to make an HTTP request.
func (surf *Surf) httpRequest(param *Request) (resp *http.Response, err error) {
	req, err := http.NewRequest(param.Method, param.Url, param.body)
	if err != nil {
		return nil, err
	}

	req.Header = param.Header

	if param.TryTimes <= 0 {
		for {
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.EnableCookie {
					l := len(UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(param.RetryPause)
				param.renewBody()
				continue
			}
			break
		}
	} else {
		for i := 0; i < param.TryTimes; i++ {
			if i != 0 {
				time.Sleep(param.RetryPause)
				param.renewBody()
			}
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.EnableCookie {
					l := len(UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", UserAgents["common"][r.Intn(l)])
				}
				continue
			}
			break
		}
	}

	return resp, err
}
