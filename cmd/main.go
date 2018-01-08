package main

import (
	"os"
	"os/signal"
	"syscall"

	brain "github.com/sundy-li/wechat_brain"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGUSR1, syscall.SIGUSR2)
	go func() {
		brain.Run("8998")
		//阻塞直至有信号传入
	}()
	<-c
	brain.Close()
}
