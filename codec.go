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
	"errors"
	"net/url"
	"unsafe"
)

type Codec struct {
	ContentType string
	Write       func(v interface{}) ([]byte, error)
	Read        func([]byte, interface{}) error
}

var Form = &Codec{
	ContentType: "application/x-www-form-urlencoded",
	Marshal: func(v interface{}) ([]byte, error) {
		switch value := v.(type) {
		case string:
			return *(*[]byte)(unsafe.Pointer(&value)), nil
		case url.Values:
			s = value.Encode()
			return *(*[]byte)(unsafe.Pointer(&s)), nil
		default:
			return []byte{}, errors.New("v needs to be string or url.Values type")
		}
	},
	Unmarshal: func(data []byte, v interface{}) error {
		if len(data) == 0 {
			return nil
		}
		switch value := v.(type) {
		case *string:
			*value = string(data)
			return nil
		case *url.Values:
			vs, err := url.ParseQuery(*(*string)(unsafe.Pointer(&data)))
			if err == nil {
				return err
			}
			*value = vs
			return nil
		default:
			return errors.New("v needs to be *url.Values type")
		}
	},
}

func parsePostForm(r *Request) (vs url.Values, err error) {
	if r.Body == nil {
		err = errors.New("missing form body")
		return
	}
	ct := r.Header.Get("Content-Type")
	// RFC 2616, section 7.2.1 - empty type
	//   SHOULD be treated as application/octet-stream
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err = mime.ParseMediaType(ct)
	switch {
	case ct == "application/x-www-form-urlencoded":
		var reader io.Reader = r.Body
		maxFormSize := int64(1<<63 - 1)
		if _, ok := r.Body.(*maxBytesReader); !ok {
			maxFormSize = int64(10 << 20) // 10 MB is a lot of text.
			reader = io.LimitReader(r.Body, maxFormSize+1)
		}
		b, e := ioutil.ReadAll(reader)
		if e != nil {
			if err == nil {
				err = e
			}
			break
		}
		if int64(len(b)) > maxFormSize {
			err = errors.New("http: POST too large")
			return
		}
		vs, e = url.ParseQuery(string(b))
		if err == nil {
			err = e
		}
	case ct == "multipart/form-data":
		// handled by ParseMultipartForm (which is calling us, or should be)
		// TODO(bradfitz): there are too many possible
		// orders to call too many functions here.
		// Clean this up and write more tests.
		// request_test.go contains the start of this,
		// in TestParseMultipartFormOrder and others.
	}
	return
}

var JOSN = &Codec{
	ContentType: "text/json",
	Marshal:     json.Marshal,
	Unmarshal:   json.Unmarshal,
}
