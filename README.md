# wechat_brain
小程序头脑王者辅助工具，上万题库。


## 注意
本工具仅供辅助娱乐。

## 使用原理
本工具运行在PC端,本质是一个http/https代理服务,对头脑王者的接口请求进行截获,主要作用有

- 将题目和返回的标准答案存储于本地数据库中(questions.data文件)
- 未匹配到标准答案情况下,自动请求搜索引擎,注解形式返回最佳概率结果

## 使用步骤：
本工具必须结合PC和手机共同使用,PC和手机须在同一个网络下

#### 以下为PC电脑操作步骤

- 运行主程序。运行方法（三选一)
	1. 方法一: 在[release](https://github.com/sundy-li/wechat_brain/releases)页面下载对应的操作系统执行文件, 解压后, 将最新版本的[questions.data](https://github.com/sundy-li/wechat_brain/blob/master/questions.data) 文件下载到同一个目录, 然后运行brain文件即可,命令行输入`./brain` 
	2. 方法二: 安装go(>=1.8)环境后, clone本repo源码到对应`$GOPATH/src/github.com/sundy-li/`下, 进入源码目录后,执行 `go run cmd/main.go`。
	3. 方法三: 使用docker命令运行：`docker build . -t wechat_brain  &&  docker run -p 8998:8998 --name my_wechat_brain -d wechat_brain`

- 新版本(version >= v0.18)加入了三种模式, 大家根据自己的需求选择模式运行
	1. 模式一: 默认模式, 修改了服务端返回的数据, 更加友好地提示正确答案, 运行方式如上所述: `./brain` 或者源码下执行 `go run cmd/main.go`
	2. 模式二: 隐身模式, 严格返回原始数据, 该模式可以防止作弊检测(客户端提交返回题目和服务端对比,模式一很容易被侦测出使用了作弊, 模式二避免了这类检测), 但该模式的缺点是降低了用户的体验,题目答案的提示只能在PC电脑上显示, 运行方式如上所述 `./brain -m 1` 或者源码下执行 `go run cmd/main.go -m 1`
	3. 模式三：自动模式 ** 注意此模式不同手机点击可能不稳定, 谨慎使用 ** 安卓机的自动刷题模式，需要将手机连接到电脑，并安装adb，且需要在开发者模式中打开usb调试，使用前请根据自身手机分辨率，调整spider文件clickProcess中的相应参数：手机屏幕中心x坐标，第一个选项中心y坐标，排位列表中最后一项中心y坐标。运行方式如上所述 `./brain -a 1 -m 1` 或者源码下执行 `go run cmd/main.go -a 1 -m 1`

#### 以下为手机安装步骤

- 设置手机代理。手机连接wifi后进行代理设置，代理IP为个人pc的内网ip地址,以及端口为8998,移动网络下可通过设置新建APN并在其中设置代理的方式实现。如：
<div align="center">    
 <img src="./docs/3.jpeg" width = "300" alt="配置代理" align=center />
</div> 

- 安装证书。代理运行成功后,手机浏览器访问 `abc.com`安装证书,ios记得要信任证书 (或者将 `certs/goproxy.crt`传到手机, 点击安装证书), 很多朋友会卡在安装证书这一步骤, 不同手机会有不同的安装方式,建议大家多搜索下自己机型如何安装证书

- 打开微信并启动头脑王者小程序。
- 正确的答案将在小程序的选项中以【标准答案】或【数字】字样。如：  
<div align="center">
 <img src="./docs/2.jpg" width = "300" alt="自动提示标准答案" align=center />
 <img src="./docs/1.jpg" width = "300" alt="自动估算最可能的答案" align=center />
</div>

## 问题
- 感谢@HsiangHo, @milkmeowo 的贡献,修复了ios代理问题,更新新版本后,最好重新安装证书,重启微信进程
  ~~~ios端由于goproxy无法代理websocket问题,暂时无法使用,希望大家可以来完善这个问题,见[这个issue](https://github.com/sundy-li/wechat_brain/issues/18)~~~

## 合并题库
- 请将questions.data文件压缩为zip文件后提交到[这里](https://github.com/sundy-li/wechat_brain/issues/17), 题库将会定期合并更新。

## 轻松上王者效果图

<div align="center">    
 <img src="./docs/4.jpeg" width = "300" alt="自动提示标准答案" align=center />
</div>

