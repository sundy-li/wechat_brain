package wechat_brain

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
)

var (
	memoryDb       *bolt.DB
	QuestionBucket = "Question"
)

func init() {
	var err error
	memoryDb, err = bolt.Open("questions.data", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	memoryDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(QuestionBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}

func StoreQuestion(question *Question) error {
	if question.CalData.TrueAnswer != "" {
		return memoryDb.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(QuestionBucket))
			v := NewQuestionCols(question.CalData.TrueAnswer)
			err := b.Put([]byte(question.Data.Quiz), v.GetData())
			return err
		})
	}
	return nil
}

func FetchQuestion(question *Question) (str string) {
	memoryDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(QuestionBucket))
		v := b.Get([]byte(question.Data.Quiz))
		if len(v) == 0 {
			return nil
		}
		q := DecodeQuestionCols(v, time.Now().Unix())
		str = q.Answer
		return nil
	})
	return
}

func FetchQuestionTime(quiz string) (res int64) {
	memoryDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(QuestionBucket))
		v := b.Get([]byte(quiz))
		if len(v) == 0 {
			res = -1
			return nil
		}
		q := DecodeQuestionCols(v, time.Now().Unix())
		res = q.Update
		return nil
	})
	return
}

func ShowAllQuestions() {
	var kv = map[string]string{}
	memoryDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(QuestionBucket))
		c := b.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			fmt.Printf("key=%s, value=%s\n", k, v)
			kv[string(k)] = string(v)
		}
		return nil
	})

}

func CountQuestions() int {
	var i int
	memoryDb.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(QuestionBucket))
		c := b.Cursor()

		for k, _ := c.First(); k != nil; k, _ = c.Next() {
			i++
		}
		return nil
	})
	return i
}

func MergeQuestions(fs ...string) {
	var i int
	for _, f := range fs {
		thirdDb, err := bolt.Open(f, 0600, nil)
		if err != nil {
			log.Println("error in merge file db "+f, err.Error())
			continue
		}
		defer thirdDb.Close()
		thirdDb.View(func(thirdTx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := thirdTx.Bucket([]byte(QuestionBucket))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				memoryDb.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(QuestionBucket))
					//三方包的时间
					q := DecodeQuestionCols(v, 0)
					//数据库中的时间
					if q.Update > FetchQuestionTime(string(k)) {
						i++
						b.Put(k, q.GetData())
					}
					return nil
				})
			}
			log.Println("merged file", f)
			return nil
		})
	}
	log.Println("merged", i, "questions")
}

type QuestionCols struct {
	Answer string `json:"a"`
	Update int64  `json:"ts"`
}

func NewQuestionCols(answer string) *QuestionCols {
	return &QuestionCols{
		Answer: answer,
		Update: time.Now().Unix(),
	}
}

func DecodeQuestionCols(bs []byte, update int64) *QuestionCols {
	var q = &QuestionCols{}
	err := json.Unmarshal(bs, q)
	if err == nil {
		return q
	} else {
		q = NewQuestionCols(string(bs))
		q.Update = update
	}
	return q
}

func (q *QuestionCols) GetData() []byte {
	bs, _ := json.Marshal(q)
	return bs
}

func max(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
