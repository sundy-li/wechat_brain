package conf

type Config struct {
	// 下载线程数
	DownloadThread int

	// phantomjs 执行文件路径
	Phantomjs string
	// js临时目录
	PhantomjsTemp string
	PuppeteerUrl  string
}

var (
	baseConfig = &Config{
		DownloadThread: DOWNLOAD_THREAD,
		Phantomjs:      "/usr/local/phantomjs/bin/phantomjs",
		PhantomjsTemp:  "/tmp/",
		PuppeteerUrl:   "http://localhost:8000",
	}
)

func Conf() Config {
	return *baseConfig
}

func InitConfig(config Config) {
	if config.DownloadThread != 0 {
		baseConfig.DownloadThread = config.DownloadThread
	}
	if config.Phantomjs != "" {
		baseConfig.Phantomjs = config.Phantomjs
	}
	if config.PhantomjsTemp != "" {
		baseConfig.PhantomjsTemp = config.PhantomjsTemp
	}
}

// Some default values
var (
	DOWNLOAD_THREAD = 50
)
