package main

import (
	"os"
	"os/signal"

	brain "github.com/sundy-li/wechat_brain"
)

func main() {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill)
	go func() {
		brain.Run("8998")
	}()
	<-c
	brain.Close()
}
