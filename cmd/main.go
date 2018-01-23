package main

import (
	"flag"
	"os"
	"os/signal"

	brain "github.com/sundy-li/wechat_brain"
)

var (
	mode int
)

func init() {
	flag.IntVar(&mode, "m", 0, "run mode 0 : default mode, easy to be detected of cheating; 1 : invisible mode")
	flag.Parse()
}

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		brain.Run("8998", mode)
	}()
	<-c
	brain.Close()
}
