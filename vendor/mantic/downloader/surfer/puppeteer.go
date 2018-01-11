package surfer

import (
	"net/http"
	"net/http/cookiejar"
)

type Puppeteer struct {
	nodeUrl string
}

//nodeUrl是puppeteer写的一个http服务的地址
func NewPuppeteer(nodeUrl string) Surfer {
	return &Puppeteer{
		nodeUrl: nodeUrl,
	}
}

func (p *Puppeteer) Download(req Request) (resp *http.Response, err error) {
	clinet := http.Client{}
	getUrl := p.nodeUrl + "?url=" + req.GetUrl()
	r, _ := http.NewRequest("GET", getUrl, nil)
	return clinet.Do(r)
}

func (p *Puppeteer) SetCookieJar(jar *cookiejar.Jar) {
	return
}

func (p *Puppeteer) Jar() *cookiejar.Jar {
	return nil
}
