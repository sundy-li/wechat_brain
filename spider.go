package wechat_brain

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"
	"time"
	"strconv"

	"github.com/coreos/goproxy"
)

var (
	_spider = newSpider()
	Mode    int
)

type spider struct {
	proxy *goproxy.ProxyHttpServer
}

func Run(port string, mode int) {
	Mode = mode
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
		if ctx.Req.URL.Path == `/question/bat/findQuiz` || ctx.Req.URL.Path == `/question/fight/findQuiz` {
			bs, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewReader(bs))
			handleQuestionReq(bs)
		} else if ctx.Req.URL.Path == `/question/bat/choose` || ctx.Req.URL.Path == `/question/fight/choose` {
			bs, _ := ioutil.ReadAll(req.Body)
			req.Body = ioutil.NopCloser(bytes.NewReader(bs))
			handleChooseReq(bs)
		} else if ctx.Req.URL.Host == `abc.com` {
			resp = new(http.Response)
			resp.StatusCode = 200
			resp.Header = make(http.Header)
			resp.Header.Add("Content-Disposition", "attachment; filename=ca.crt")
			resp.Header.Add("Content-Type", "application/octet-stream")
			resp.Body = ioutil.NopCloser(bytes.NewReader(goproxy.CA_CERT))
		}
		return
	}
	responseHandleFunc := func(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
		if resp == nil { return resp }
		switch ctx.Req.URL.Path {
			case "/question/bat/findQuiz":
				bs, _ := ioutil.ReadAll(resp.Body)
				bsNew,ansPos := handleQuestionResp(bs)
				resp.Body = ioutil.NopCloser(bytes.NewReader(bsNew))
				go clickProcess(ansPos) // click answer
				break
			case "/question/bat/choose":
				bs, _ := ioutil.ReadAll(resp.Body)
				resp.Body = ioutil.NopCloser(bytes.NewReader(bs))
				go handleChooseResponse(bs)
				break
			case "/question/bat/fightResult":
				go clickProcess(-1) // go to next match
				break
		}
		return resp
	}
	s.proxy.OnResponse().DoFunc(responseHandleFunc)
	s.proxy.OnRequest().DoFunc(requestHandleFunc)
}

func clickProcess(ansPos int) {
	var enableFlag = false // control flag
	var screanCenterX = 550 // center of screen
	var firstItemY = 1280 // center of first item (y)
	var qualifyingItemY = 2000 // 排位列表最后一项 y 坐标
	
	if(enableFlag) {
		if(ansPos >= 0) {
			log.Printf("【点击】正在点击选项：%d", ansPos)
			time.Sleep(time.Millisecond * 3800) //延迟
			go clickAction(screanCenterX, firstItemY + 200 * (ansPos - 1)) // process click
		}else{
			// go to next match
			log.Printf("【点击】将点击继续挑战按钮...")
			time.Sleep(time.Millisecond * 7500)
			go clickAction(screanCenterX, firstItemY + 400) // 继续挑战 按钮在第三个item处
			log.Printf("【点击】将点击排位列表底部一项，进行比赛匹配...")
			time.Sleep(time.Millisecond * 2000)
			go clickAction(screanCenterX, qualifyingItemY)
		}
	}
}

func clickAction(posX int, posY int) {
	var err error
	touchX, touchY := strconv.Itoa(posX), strconv.Itoa(posY)
	_, err = exec.Command("adb","shell", "input", "swipe", touchX, touchY,touchX, touchY).Output()
	if err != nil { log.Fatal("error: check adb connection.") }
}

func orPanic(err error) {
	if err != nil {
		panic(err)
	}
}
