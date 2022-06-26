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
	"io/ioutil"
	"net/http"
	"testing"
)

func TestSurf(t *testing.T) {
	req := &Request{
		Method:       "GET",
		Url:          "https://www.bing.com/search?q=golang",
		EnableCookie: true,
		Header: http.Header{
			"Origin": []string{"https://cn.bing.com"},
		},
	}
	resp, _ := Download(req)
	b, _ := ioutil.ReadAll(resp.Body)
	t.Logf("request:\n%#v", req)
	t.Logf("response:\n%#v\nresponse_body:\n%s", resp, b[:200])
	resp, _ = Download(req)
	b, _ = ioutil.ReadAll(resp.Body)
	t.Logf("request:\n%#v", req)
	t.Logf("response:\n%#v\nresponse_body:\n%s", resp, b[:200])
}
