package wechat_brain

import (
	"time"

	cache "github.com/patrickmn/go-cache"
)

var (
	questions *cache.Cache
)

func init() {
	// Create a cache with a default expiration time of 5 minutes, and which
	// purges expired items every 10 minutes
	questions = cache.New(5*time.Minute, 10*time.Minute)
}

func GetQuestion(roomID, quizNum string) *Question {
	key := roomID + "_" + quizNum
	if entity, ok := questions.Get(key); ok {
		return entity.(*Question)
	}
	return nil
}

func SetQuestion(question *Question) {
	key := question.CalData.RoomID + "_" + question.CalData.quizNum
	questions.Set(key, question, cache.DefaultExpiration)
}
