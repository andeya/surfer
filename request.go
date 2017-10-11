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
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// A Request represents an HTTP request received by a server
// or to be sent by a client.
//
// The field semantics differ slightly between client and server
// usage. In addition to the notes on the fields below, see the
// documentation for Request.Write and RoundTripper.
type Request struct {
	// Method specifies the HTTP method (GET, POST, PUT, etc.).
	// For client requests an empty string means GET.
	Method string

	// URL specifies either the URI being requested (for server
	// requests) or the URL to access (for client requests).
	//
	// For server requests the URL is parsed from the URI
	// supplied on the Request-Line as stored in RequestURI.  For
	// most requests, fields other than Path and RawQuery will be
	// empty. (See RFC 2616, Section 5.1.2)
	//
	// For client requests, the URL's Host specifies the server to
	// connect to, while the Request's Host field optionally
	// specifies the Host header value to send in the HTTP
	// request.
	URL *url.URL

	// The protocol version for incoming server requests.
	//
	// For client requests these fields are ignored. The HTTP
	// client code always uses either HTTP/1.1 or HTTP/2.
	// See the docs on Transport for details.
	__Proto      string // "HTTP/1.0"
	__ProtoMajor int    // 1
	__ProtoMinor int    // 0

	// Header contains the request header fields either received
	// by the server or to be sent by the client.
	//
	// If a server received a request with header lines,
	//
	//	Host: example.com
	//	accept-encoding: gzip, deflate
	//	Accept-Language: en-us
	//	fOO: Bar
	//	foo: two
	//
	// then
	//
	//	Header = map[string][]string{
	//		"Accept-Encoding": {"gzip, deflate"},
	//		"Accept-Language": {"en-us"},
	//		"Foo": {"Bar", "two"},
	//	}
	//
	// For incoming requests, the Host header is promoted to the
	// Request.Host field and removed from the Header map.
	//
	// HTTP defines that header names are case-insensitive. The
	// request parser implements this by using CanonicalHeaderKey,
	// making the first character and any characters following a
	// hyphen uppercase and the rest lowercase.
	//
	// For client requests, certain headers such as Content-Length
	// and Connection are automatically written when needed and
	// values in Header may be ignored. See the documentation
	// for the Request.Write method.
	Header Header

	// Body is the request's body.
	//
	// For client requests a nil body means the request has no
	// body, such as a GET request. The HTTP Client's Transport
	// is responsible for calling the Close method.
	//
	// For server requests the Request Body is always non-nil
	// but will return EOF immediately when no body is present.
	// The Server will close the request body. The ServeHTTP
	// Handler does not need to.
	Body io.ReadCloser

	// GetBody defines an optional func to return a new copy of
	// Body. It is used for client requests when a redirect requires
	// reading the body more than once. Use of GetBody still
	// requires setting Body.
	//
	// For server requests it is unused.
	_GetBody func() (io.ReadCloser, error)

	// ContentLength records the length of the associated content.
	// The value -1 indicates that the length is unknown.
	// Values >= 0 indicate that the given number of bytes may
	// be read from Body.
	// For client requests, a value of 0 with a non-nil Body is
	// also treated as unknown.
	ContentLength int64

	// TransferEncoding lists the transfer encodings from outermost to
	// innermost. An empty list denotes the "identity" encoding.
	// TransferEncoding can usually be ignored; chunked encoding is
	// automatically added and removed as necessary when sending and
	// receiving requests.
	TransferEncoding []string

	// Close indicates whether to close the connection after
	// replying to this request (for servers) or after sending this
	// request and reading its response (for clients).
	//
	// For server requests, the HTTP server handles this automatically
	// and this field is not needed by Handlers.
	//
	// For client requests, setting this field prevents re-use of
	// TCP connections between requests to the same hosts, as if
	// Transport.DisableKeepAlives were set.
	_Close bool

	// For server requests Host specifies the host on which the
	// URL is sought. Per RFC 2616, this is either the value of
	// the "Host" header or the host name given in the URL itself.
	// It may be of the form "host:port". For international domain
	// names, Host may be in Punycode or Unicode form. Use
	// golang.org/x/net/idna to convert it to either format if
	// needed.
	//
	// For client requests Host optionally overrides the Host
	// header to send. If empty, the Request.Write method uses
	// the value of URL.Host. Host may contain an international
	// domain name.
	_Host string

	// Form contains the parsed form data, including both the URL
	// field's query parameters and the POST or PUT form data.
	// This field is only available after ParseForm is called.
	// The HTTP client ignores Form and uses Body instead.
	__Form url.Values

	// PostForm contains the parsed form data from POST, PATCH,
	// or PUT body parameters.
	//
	// This field is only available after ParseForm is called.
	// The HTTP client ignores PostForm and uses Body instead.
	__PostForm url.Values

	// MultipartForm is the parsed multipart form, including file uploads.
	// This field is only available after ParseMultipartForm is called.
	// The HTTP client ignores MultipartForm and uses Body instead.
	__MultipartForm *multipart.Form

	// Trailer specifies additional headers that are sent after the request
	// body.
	//
	// For server requests the Trailer map initially contains only the
	// trailer keys, with nil values. (The client declares which trailers it
	// will later send.)  While the handler is reading from Body, it must
	// not reference Trailer. After reading from Body returns EOF, Trailer
	// can be read again and will contain non-nil values, if they were sent
	// by the client.
	//
	// For client requests Trailer must be initialized to a map containing
	// the trailer keys to later send. The values may be nil or their final
	// values. The ContentLength must be 0 or -1, to send a chunked request.
	// After the HTTP request is sent the map values can be updated while
	// the request body is read. Once the body returns EOF, the caller must
	// not mutate Trailer.
	//
	// Few HTTP clients, servers, or proxies support HTTP trailers.
	__Trailer Header

	// RemoteAddr allows HTTP servers and other software to record
	// the network address that sent the request, usually for
	// logging. This field is not filled in by ReadRequest and
	// has no defined format. The HTTP server in this package
	// sets RemoteAddr to an "IP:port" address before invoking a
	// handler.
	// This field is ignored by the HTTP client.
	__RemoteAddr string

	// RequestURI is the unmodified Request-URI of the
	// Request-Line (RFC 2616, Section 5.1) as sent by the client
	// to a server. Usually the URL field should be used instead.
	// It is an error to set this field in an HTTP client request.
	RequestURI string

	// TLS allows HTTP servers and other software to record
	// information about the TLS connection on which the request
	// was received. This field is not filled in by ReadRequest.
	// The HTTP server in this package sets the field for
	// TLS-enabled connections before invoking a handler;
	// otherwise it leaves the field nil.
	// This field is ignored by the HTTP client.
	__TLS *tls.ConnectionState

	// Cancel is an optional channel whose closure indicates that the client
	// request should be regarded as canceled. Not all implementations of
	// RoundTripper may support Cancel.
	//
	// For server requests, this field is not applicable.
	//
	// Deprecated: Use the Context and WithContext methods
	// instead. If a Request's Cancel field and context are both
	// set, it is undefined whether Cancel is respected.
	___Cancel <-chan struct{}

	// Response is the redirect response which caused this request
	// to be created. This field is only populated during client
	// redirects.
	Response *Response

	// ctx is either the client or server context. It should only
	// be modified via copying the whole Request using WithContext.
	// It is unexported to prevent people from using Context wrong
	// and mutating the contexts held by callers of the same request.
	ctx context.Context

	BodyCodec *Codec
}

// constant
const (
	SurfID             = 0               // Surf下载器标识符
	PhomtomJsID        = 1               // PhomtomJs下载器标识符
	DefaultMethod      = "GET"           // 默认请求方法
	DefaultDialTimeout = 2 * time.Minute // 默认请求服务器超时
	DefaultConnTimeout = 2 * time.Minute // 默认下载超时
	DefaultTryTimes    = 3               // 默认最大下载次数
	DefaultRetryPause  = 2 * time.Second // 默认重新下载前停顿时长
)

// Request contains the necessary prerequisite information.
type Request struct {
	// url (必须填写)
	Url string
	url *url.URL
	// GET POST HEAD (默认为GET)
	Method string
	// http header
	Header http.Header
	// 是否使用cookies，在Spider的EnableCookie设置
	EnableCookie bool
	// request body interface
	Body body
	body io.Reader
	// dial tcp: i/o timeout
	DialTimeout time.Duration
	// WSARecv tcp: i/o timeout
	ConnTimeout time.Duration
	// the max times of download
	TryTimes int
	// how long pause when retry
	RetryPause time.Duration
	// max redirect times
	// when RedirectTimes equal 0, redirect times is ∞
	// when RedirectTimes less than 0, redirect times is 0
	RedirectTimes int
	// the download ProxyHost
	Proxy string
	proxy *url.URL
	// 指定下载器ID
	// 0为Surf高并发下载器，各种控制功能齐全
	// 1为PhantomJS下载器，特点破防力强，速度慢，低并发
	DownloaderID int
	client       *http.Client
}

func (r *Request) prepare() error {
	var err error
	r.url, err = UrlEncode(r.Url)
	if err != nil {
		return err
	}
	r.Url = r.url.String()
	if r.Proxy != "" {
		if r.proxy, err = url.Parse(r.Proxy); err != nil {
			return err
		}
	}
	if r.DialTimeout < 0 {
		r.DialTimeout = 0
	} else if r.DialTimeout == 0 {
		r.DialTimeout = DefaultDialTimeout
	}

	if r.ConnTimeout < 0 {
		r.ConnTimeout = 0
	} else if r.ConnTimeout == 0 {
		r.ConnTimeout = DefaultConnTimeout
	}

	if r.TryTimes == 0 {
		r.TryTimes = DefaultTryTimes
	}

	if r.RetryPause <= 0 {
		r.RetryPause = DefaultRetryPause
	}

	if r.DownloaderID != PhomtomJsID {
		r.DownloaderID = SurfID
	}

	if r.Header == nil {
		r.Header = make(http.Header)
	}
	var commonUserAgentIndex int
	if !r.EnableCookie {
		commonUserAgentIndex = rand.Intn(len(UserAgents["common"]))
		r.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	} else if len(r.Header["User-Agent"]) == 0 {
		r.Header.Set("User-Agent", UserAgents["common"][commonUserAgentIndex])
	}
	if len(r.Method) == 0 {
		r.Method = DefaultMethod
	} else {
		r.Method = strings.ToUpper(r.Method)
	}
	r.body = nil
	if r.Body != nil {
		return r.Body.SetBody(r)
	}
	return nil
}

// ReadBody returns body bytes
func (r *Request) ReadBody() (b []byte, err error) {
	if r.url == nil {
		r.prepare()
	}
	if r.body != nil {
		b, err = ioutil.ReadAll(r.body)
	}
	return b, err
}

// 回写Request内容
func (r *Request) writeback(resp *http.Response) *http.Response {
	if resp == nil {
		resp = new(http.Response)
		resp.Request = new(http.Request)
	} else if resp.Request == nil {
		resp.Request = new(http.Request)
	}

	resp.Request.Method = r.Method
	resp.Request.Header = r.Header
	resp.Request.Host = r.url.Host

	// reset url
	r.url = nil

	return resp
}

// checkRedirect is used as the value to http.Client.CheckRedirect
// when redirectTimes equal 0, redirect times is ∞
// when redirectTimes less than 0, not allow redirects
func (r *Request) checkRedirect(req *http.Request, via []*http.Request) error {
	if r.RedirectTimes == 0 {
		return nil
	}
	if len(via) >= r.RedirectTimes {
		if r.RedirectTimes < 0 {
			return fmt.Errorf("not allow redirects.")
		}
		return fmt.Errorf("stopped after %v redirects.", r.RedirectTimes)
	}
	return nil
}
