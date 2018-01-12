# wechat_brain
小程序王者头脑辅助工具


## 注意
本工具仅供辅助娱乐

## 如何使用
	
- 手机浏览器访问 `abc.com`安装证书 (或者将 `certs/goproxy.crt`传到手机, 点击安装证书)
- 运行方法一: 在[release](https://github.com/sundy-li/wechat_brain/releases)页面下载对应的操作系统执行文件, 压缩后, 将`questions.data`文件下载到同一个目录
- 运行方法二: 安装go环境后, clone本repo源码到对应`$GOPATH/src/github.com/sundy-li/`下, 进入源码目录后,执行 `go run cmd/main.go`
- 手机wifi设置代理为pc的ip和端口,启动小程序王者头脑


## 问题

- ios端由于goproxy无法代理websocket问题,暂时无法使用,希望大家可以来完善这个问题,见[这个issue](https://github.com/sundy-li/wechat_brain/issues/18)


## 合并题库
- 提交到[这里](https://github.com/sundy-li/wechat_brain/issues/17), 程序会定期合并

## 运行示例

<div align="center">    
 <img src="./docs/3.jpeg" width = "300" alt="配置代理" align=center />
 <img src="./docs/2.jpg" width = "300" alt="自动提示标准答案" align=center />
 <img src="./docs/1.jpg" width = "300" alt="自动估算最可能的答案" align=center />
 <img src="./docs/4.jpeg" width = "300" alt="自动提示标准答案" align=center />
</div>

- 注意,第一个图片的host和端口是我的云主机ip端口,你添加证书后可以用这个ip端口测试下, 测试通过后在pc跑起来你的服务后, 将此改为你的电脑ip和端口