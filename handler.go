package wechat_brain

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"time"
)

var (
	roomID string
)

func handleQuestionReq(bs []byte) {
	values, _ := url.ParseQuery(string(bs))
	roomID = values.Get("roomID")
}

//根据题目返回,进行答案搜索
func handleQuestionResp(bs []byte) (bsNew []byte) {
	bsNew = bs
	question := &Question{}
	json.Unmarshal(bs, question)
	question.CalData.RoomID = roomID
	question.CalData.quizNum = strconv.Itoa(question.Data.Num)

	//Get the answer from the db
	answer := FetchQuestion(question)
	var ret map[string]int
	if answer == "" {
		tx := time.Now()
		ret = GetFromBaidu(question.Data.Quiz, question.Data.Options)
		tx2 := time.Now()
		log.Printf("Cost time %d ms\n", tx2.Sub(tx).Nanoseconds()/1e6)
	}
	question.CalData.TrueAnswer = answer
	question.CalData.Answer = answer
	SetQuestion(question)

	respQuestion := &Question{}
	json.Unmarshal(bs, respQuestion)
	if question.CalData.TrueAnswer != "" {
		for i, option := range respQuestion.Data.Options {
			if option == question.CalData.TrueAnswer {
				respQuestion.Data.Options[i] = option + "[标准答案]"
				break
			}
		}
	} else {
		for i, option := range respQuestion.Data.Options {
			if ret[option] > 0 {
				respQuestion.Data.Options[i] = option + "[" + strconv.Itoa(ret[option]) + "]"
			}
		}
	}
	bsNew, _ = json.Marshal(respQuestion)
	log.Printf("Response findQuiz%v\n", string(bsNew))
	return bsNew
}

//hijack 提交请求
func handleChooseReq(bs []byte) (newBs []byte) {
	newBs = bs
	values, _ := url.ParseQuery(string(bs))
	quizNum := values.Get("quizNum")
	question := GetQuestion(roomID, quizNum)
	if question == nil {
		return
	}
	// var idx = -1
	// for i, option := range question.Data.Options {
	// 	if question.CalData.Answer == option {
	// 		idx = i + 1 //TODO ?
	// 		break
	// 	}
	// }
	// if idx != -1 {
	// 	values.Set("option", strconv.Itoa(idx))
	// 	newBs = []byte(values.Encode())
	// }
	return
}

func handleChooseResponse(bs []byte) {
	chooseResp := &ChooseResp{}
	json.Unmarshal(bs, chooseResp)

	log.Println("response choose", roomID, chooseResp.Data.Num, string(bs))
	question := GetQuestion(roomID, strconv.Itoa(chooseResp.Data.Num))
	if question == nil {
		log.Println("error in get question", chooseResp.Data.RoomID, chooseResp.Data.Num)
		return
	}
	question.CalData.TrueAnswer = question.Data.Options[chooseResp.Data.Answer-1]
	log.Printf("Saving %s , %s", question.Data.Quiz, question.CalData.TrueAnswer)
	StoreQuestion(question)
}

type Question struct {
	Data struct {
		Quiz        string   `json:"quiz"`
		Options     []string `json:"options"`
		Num         int      `json:"num"`
		School      string   `json:"school"`
		Type        string   `json:"type"`
		Contributor string   `json:"contributor"`
		EndTime     int      `json:"endTime"`
		CurTime     int      `json:"curTime"`
	} `json:"data"`
	Errcode int `json:"errcode"`

	CalData struct {
		RoomID     string
		quizNum    string
		Answer     string
		TrueAnswer string
	} `json:"-"`
}

type ChooseResp struct {
	Data struct {
		UID         int  `json:"uid"`
		Num         int  `json:"num"`
		Answer      int  `json:"answer"`
		Option      int  `json:"option"`
		Yes         bool `json:"yes"`
		Score       int  `json:"score"`
		TotalScore  int  `json:"totalScore"`
		RowNum      int  `json:"rowNum"`
		RowMult     int  `json:"rowMult"`
		CostTime    int  `json:"costTime"`
		RoomID      int  `json:"roomId"`
		EnemyScore  int  `json:"enemyScore"`
		EnemyAnswer int  `json:"enemyAnswer"`
	} `json:"data"`
	Errcode int `json:"errcode"`
}

//roomID=476376430&quizNum=4&option=4&uid=26394007&t=1515326786076&sign=3592b9d28d045f3465206b4147ea872b
