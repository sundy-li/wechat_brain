package main

import (
	"flag"
	"log"
	"strings"

	brain "github.com/sundy-li/wechat_brain"
)

var (
	action string
	fs     string
)

func init() {
	flag.StringVar(&action, "a", "show", "action value -> show | merge")
	flag.StringVar(&fs, "fs", "", "merge data files")
	flag.Parse()
}

func main() {
	if action == "merge" {
		files := strings.Split(fs, " ")
		if len(files) < 1 {
			log.Println("empty files")
			return
		}
		brain.MergeQuestions(files...)
	}
	total := brain.CountQuestions()
	log.Println("total questions =>", total)
}
