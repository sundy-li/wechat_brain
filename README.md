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

- 运行主程序。运行方法（二选一)
	1. 方法一: 在[release](https://github.com/sundy-li/wechat_brain/releases)页面下载对应的操作系统执行文件, 解压后, 将最新版本的[questions.data](https://github.com/sundy-li/wechat_brain/blob/master/questions.data) 文件下载到同一个目录, 然后运行brain文件即可,命令行输入`./brain` 
	2. 方法二: 安装go(>=1.8)环境后, clone本repo源码到对应`$GOPATH/src/github.com/sundy-li/`下, 进入源码目录后,执行 `go run cmd/main.go`。

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

