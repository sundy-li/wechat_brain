package wechat_brain

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/elazarl/goproxy"
)

var (
	_spider = newSpider()
)

type spider struct {
	proxy *goproxy.ProxyHttpServer
}

func Run(port string) {
	_spider.Init()
	_spider.Run(port)
}

func Close() {
	memoryDb.Close()
}

func newSpider() *spider {
	sp := &spider{}
	sp.proxy = goproxy.NewProxyHttpServer()
	sp.proxy.OnRequest().HandleConnect(goproxy.AlwaysMitm)
	return sp
}

func (s *spider) Run(port string) {
	log.Println("server will at port:" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.proxy))
}

func (s *spider) Init() {
	requestHandleFunc := func(request *http.Request, ctx *goproxy.ProxyCtx) (req *http.Request, resp *http.Response) {
		req = request
		if ctx.Req.URL.Path == `/question/fight/findQuiz` {
			bs, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewReader(bs))
			handleQuestionReq(bs)
		} else if ctx.Req.URL.Path == `/question/fight/choose` {
			bs, _ := ioutil.ReadAll(req.Body)
			log.Println("request choose==>", string(bs))
			req.Body = ioutil.NopCloser(bytes.NewReader(bs))
			handleChooseReq(bs)
		}
		return
	}
	responseHandleFunc := func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil {
			return resp
		}
		if ctx.Req.URL.Path == `/question/fight/findQuiz` {
			bs, _ := ioutil.ReadAll(resp.Body)
			bsNew := handleQuestionResp(bs)
			resp.Body = ioutil.NopCloser(bytes.NewReader(bsNew))
			log.Println("response findQuiz==>", string(bs))

		} else if ctx.Req.URL.Path == `/question/fight/choose` {
			bs, _ := ioutil.ReadAll(resp.Body)
			resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
			go handleChooseResponse(bs)
		}
		return resp
	}
	s.proxy.OnResponse().DoFunc(responseHandleFunc)
	s.proxy.OnRequest().DoFunc(requestHandleFunc)
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
