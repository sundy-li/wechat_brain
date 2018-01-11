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
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// 基于Phantomjs的下载器实现，作为surfer的补充
// 效率较surfer会慢很多，但是因为模拟浏览器，破防性更好
// 支持UserAgent/TryTimes/RetryPause/自定义js
type (
	Phantom struct {
		PhantomjsFile string            //Phantomjs完整文件名
		TempJsDir     string            //临时js存放目录
		jsFileMap     map[string]string //已存在的js文件
		cookiejar     *cookiejar.Jar
	}
	Response struct {
		Cookie string
		Body   string
	}
)

func NewPhantom(phantomjsFile, tempJsDir string) Surfer {
	phantom := &Phantom{
		PhantomjsFile: phantomjsFile,
		TempJsDir:     tempJsDir,
		jsFileMap:     make(map[string]string),
	}
	phantom.cookiejar, _ = cookiejar.New(nil)
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
	phantom.createJsFile("get", getJs)
	phantom.createJsFile("post", postJs)
	phantom.createJsFile("url", urlJs)
	return phantom
}

var reg = regexp.MustCompile(`(\{"Cookie":.+?})`)

// 实现surfer下载器接口
func (self *Phantom) Download(req Request) (resp *http.Response, err error) {
	var encoding = "utf-8"
	if _, params, err := mime.ParseMediaType(req.GetHeader().Get("Content-Type")); err == nil {
		if cs, ok := params["charset"]; ok {
			encoding = strings.ToLower(strings.TrimSpace(cs))
		}
	}
	req.GetHeader().Del("Content-Type")

	param, err := NewParam(req)
	if err != nil {
		return nil, err
	}
	resp = param.writeback(resp)
	var args []string
	switch req.GetMethod() {
	case "GET":
		args = []string{
			self.jsFileMap["get"],
			req.GetUrl(),
			self.getCookie(param),
			encoding,
			param.header.Get("User-Agent"),
			param.header.Get("Referer"),
		}
	case "POST", "POST-M":
		args = []string{
			self.jsFileMap["post"],
			req.GetUrl(),
			self.getCookie(param),
			encoding,
			param.header.Get("User-Agent"),
			param.header.Get("Referer"),
			req.GetPostData(),
		}
	case "URL":
		args = []string{
			self.jsFileMap["url"],
			req.GetUrl(),
			self.getCookie(param),
			encoding,
			param.header.Get("User-Agent"),
			strconv.Itoa(int(param.connTimeout) / 1e6),
			req.GetPostData(),
		}
	}
	for i := 0; i < param.tryTimes; i++ {
		cmd := exec.Command(self.PhantomjsFile, args...)
		bs, err := cmd.CombinedOutput()
		if cmd.Process != nil {
			cmd.Process.Kill()
		}
		if err != nil {
			time.Sleep(param.retryPause)
			continue
		}

		if req.GetMethod() == "URL" {
			match := reg.FindAllStringSubmatch(string(bs), -1)
			if len(match) < 1 {
				time.Sleep(param.retryPause)
				continue
			}
			bs = []byte(match[0][1])
		}

		retResp := Response{}
		err = json.Unmarshal(bs, &retResp)
		if err != nil {
			time.Sleep(param.retryPause)
			continue
		}
		resp.Header = param.header
		cookies := strings.Split(strings.TrimSpace(retResp.Cookie), "; ")
		for _, c := range cookies {
			resp.Header.Add("Set-Cookie", c)
		}
		self.cookiejar.SetCookies(param.url, resp.Cookies())
		resp.Body = ioutil.NopCloser(strings.NewReader(retResp.Body))
		break
	}

	if err == nil {
		resp.StatusCode = http.StatusOK
		resp.Status = http.StatusText(http.StatusOK)
	} else {
		resp.StatusCode = http.StatusBadGateway
		resp.Status = http.StatusText(http.StatusBadGateway)
	}
	return
}

func (self *Phantom) SetCookieJar(jar *cookiejar.Jar) {
	self.cookiejar = jar
}

//销毁js临时文件
func (self *Phantom) DestroyJsFiles() {
	p, _ := filepath.Split(self.TempJsDir)
	if p == "" {
		return
	}
	for _, filename := range self.jsFileMap {
		os.Remove(filename)
	}
	if len(WalkDir(p)) == 1 {
		os.Remove(p)
	}
}

func (self *Phantom) getCookie(param *Param) string {
	cookie := param.header.Get("cookie")
	cs := self.cookiejar.Cookies(param.url)
	for _, c := range cs {
		cookie += c.Name + "=" + c.Value + ";"
	}
	return cookie
}

func (self *Phantom) Jar() *cookiejar.Jar {
	return self.cookiejar
}

func (self *Phantom) createJsFile(fileName, jsCode string) {
	fullFileName := filepath.Join(self.TempJsDir, fileName)
	// 创建并写入文件
	f, _ := os.Create(fullFileName)
	f.Write([]byte(jsCode))
	f.Close()
	self.jsFileMap[fileName] = fullFileName
}

/*
* GET method
* system.args[0] == get.js
* system.args[1] == url
* system.args[2] == cookie
* system.args[3] == pageEncode
* system.args[4] == userAgent
* system.args[5] == referer
 */

const getJs string = `
var system = require('system');
var page = require('webpage').create();
var url = system.args[1];
var cookie = system.args[2];
var pageEncode = system.args[3];
var userAgent = system.args[4];
var referer = system.args[5];


if(pageEncode){
	phantom.outputEncoding = pageEncode;
}

var headers =  {};
if(userAgent){
	headers['User-Agent'] = userAgent
}
if(referer){
	headers['Referer'] = referer
}
page.customHeaders = headers;
 

page.onResourceRequested = function(requestData, request) {
	if(cookie){
    	request.setHeader('Cookie', cookie)
    }
};

page.open(url, function(status) {
    if (status !== 'success') {
        console.log('Unable to access network');
        phantom.exit();
    }  
});

page.onLoadFinished = function(status) {
    var cookie = page.evaluate(function(s) {
            return document.cookie;
        });
    var resp = {
            "Cookie": cookie,
        	"Body": page.content
    };
    console.log(JSON.stringify(resp));
    phantom.exit();
};
`

/*
* POST method
* system.args[0] == post.js
* system.args[1] == url
* system.args[2] == cookie
* system.args[3] == pageEncode
* system.args[4] == userAgent
* system.args[5] == referer
* system.args[6] == postdata
 */
const postJs string = `
var system = require('system');
var page = require('webpage').create();
var url = system.args[1];
var cookie = system.args[2];
var pageEncode = system.args[3];
var userAgent = system.args[4];
var referer = system.args[5];
var postdata = system.args[6];

if(pageEncode){
	phantom.outputEncoding = pageEncode;
}

var headers =  {};
if(userAgent){
	headers['User-Agent'] = userAgent
}
if(referer){
	headers['Referer'] = referer
}
page.customHeaders = headers;
 

page.onResourceRequested = function(requestData, request) {
	if(cookie){
    	request.setHeader('Cookie', cookie)
    }
};

page.open(url, 'post', postdata, function(status) {
     if (status !== 'success') {
        console.log('Unable to access network');
        phantom.exit();
     }
});

page.onLoadFinished = function(status) {
    var cookie = page.evaluate(function(s) {
            return document.cookie;
    });
    var resp = {
            "Cookie": cookie,
            "Body": page.content
    };
    console.log(JSON.stringify(resp));
    phantom.exit();
});
`

/*
* url method
* system.args[0] == url.js
* system.args[1] == url
* system.args[2] == cookie
* system.args[3] == pageEncode
* system.args[4] == userAgent
* system.args[5] == timeout
* system.args[6] == url
 */
const urlJs string = `
var system = require('system');
var page = require('webpage').create();
var url = system.args[1];
var cookie = system.args[2];
var pageEncode = system.args[3];
var userAgent = system.args[4];
var timeout = system.args[5];
var uri = system.args[6];

if(pageEncode){
	phantom.outputEncoding = pageEncode;
}

var headers =  {};
if(userAgent){
	headers['User-Agent'] = userAgent
}
if(referer){
	headers['Referer'] = referer
}
page.customHeaders = headers;


page.onResourceRequested = function(requestData, request) {
	if(cookie){
    	request.setHeader('Cookie', cookie)
    }
};

page.onResourceReceived = function (res) {
	reg = new RegExp(uri, "g")
	if (res.url.match(reg)) {
		var cookie = page.evaluate(function(s) {
            return document.cookie;
        });
        var resp = {
            "Cookie": cookie,
            "Body": res.url
        };
        console.log(JSON.stringify(resp));
        phantom.exit();
	}
};

page.open(url, function(status) {
    if (status !== 'success') {
        //console.log('Unable to access network');
    }
    window.setTimeout(function () {
            phantom.exit(1);
        }, timeout);
});
`
