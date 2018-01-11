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

	"mantic/cache/dns"
	"mantic/downloader/surfer/agent"
)

// Default is the default Download implementation.
type Surf struct {
	cookieJar *cookiejar.Jar
	resolver  *dns.Resolver
}

func New() Surfer {
	s := new(Surf)
	s.cookieJar, _ = cookiejar.New(nil)
	s.resolver = dns.New(15 * time.Minute)
	return s
}

func (self *Surf) Download(req Request) (resp *http.Response, err error) {
	param, err := NewParam(req)
	if err != nil {
		return nil, err
	}
	param.client = self.buildClient(param)
	resp, err = self.httpRequest(param)

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

	resp = param.writeback(resp)
	return
}

func (self *Surf) SetCookieJar(jar *cookiejar.Jar) {
	self.cookieJar = jar
}

func (self *Surf) Jar() *cookiejar.Jar {
	return self.cookieJar
}

var (
	timeout, keepalive, idle               = 30, 30, 100
	TIMEOUT                                = time.Duration(timeout) * time.Second
	resolver                 *dns.Resolver = dns.New(15 * time.Minute)
	defaultTransport                       = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: func(network string, address string) (net.Conn, error) {
			separator := strings.LastIndex(address, ":")
			ip, _ := resolver.FetchOneString(address[:separator])
			c, err := net.DialTimeout(network, ip+address[separator:], TIMEOUT)
			if err != nil {
				return nil, err
			}
			if TIMEOUT > 0 {
				c.SetDeadline(time.Now().Add(TIMEOUT))
			}
			return c, nil
		},
		ResponseHeaderTimeout: TIMEOUT,
		TLSHandshakeTimeout:   TIMEOUT,
		MaxIdleConnsPerHost:   idle,
	}
)

func (s *Surf) newTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Dial: (&net.Dialer{
			Timeout:   TIMEOUT,
			KeepAlive: time.Duration(keepalive) * time.Second,
		}).Dial,
		ResponseHeaderTimeout: TIMEOUT,
		TLSHandshakeTimeout:   TIMEOUT,
		MaxIdleConnsPerHost:   idle,
	}
}

func (s *Surf) getDefaultTransport() *http.Transport {
	return defaultTransport
}

// buildClient creates, configures, and returns a *http.Client type.
func (self *Surf) buildClient(param *Param) *http.Client {
	client := &http.Client{
		CheckRedirect: param.checkRedirect,
	}

	if param.enableCookie {
		client.Jar = self.cookieJar
	}

	// transport := &http.Transport{
	// 	Dial: func(network, address string) (net.Conn, error) {
	// 		separator := strings.LastIndex(address, ":")
	// 		ip, _ := self.resolver.FetchOneString(address[:separator])
	// 		c, err := net.DialTimeout(network, ip+address[separator:], param.dialTimeout)
	// 		if err != nil {
	// 			return nil, err
	// 		}
	// 		if param.connTimeout > 0 {
	// 			c.SetDeadline(time.Now().Add(param.connTimeout))
	// 		}
	// 		return c, nil
	// 	},
	// }

	var transport *http.Transport
	if param.proxy != nil {
		transport = self.newTransport()
		transport.Proxy = http.ProxyURL(param.proxy)
	} else {
		transport = self.getDefaultTransport()
	}

	if strings.ToLower(param.url.Scheme) == "https" {
		transport.TLSClientConfig = &tls.Config{RootCAs: nil, InsecureSkipVerify: true}
		transport.DisableCompression = true
	}
	client.Transport = transport
	return client
}

// send uses the given *http.Request to make an HTTP request.
func (self *Surf) httpRequest(param *Param) (resp *http.Response, err error) {
	req, err := http.NewRequest(param.method, param.url.String(), param.body)
	if err != nil {
		return nil, err
	}

	req.Header = param.header
	if param.tryTimes <= 0 {
		for {
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.enableCookie {
					l := len(agent.UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", agent.UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(param.retryPause)
				continue
			}
			break
		}
	} else {
		for i := 0; i < param.tryTimes; i++ {
			resp, err = param.client.Do(req)
			if err != nil {
				if !param.enableCookie {
					l := len(agent.UserAgents["common"])
					r := rand.New(rand.NewSource(time.Now().UnixNano()))
					req.Header.Set("User-Agent", agent.UserAgents["common"][r.Intn(l)])
				}
				time.Sleep(param.retryPause)
				continue
			}
			break
		}
	}
	if resp != nil {
		self.cookieJar.SetCookies(param.url, resp.Cookies())
	}
	return resp, err
}
