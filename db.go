package wechat_brain

import (
	"fmt"
	"log"

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
		_, err := tx.CreateBucket([]byte(QuestionBucket))
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
			err := b.Put([]byte(question.Data.Quiz), []byte(question.CalData.TrueAnswer))
			return err
		})
	}
	return nil
}

func FetchQuestion(question *Question) (str string) {
	memoryDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(QuestionBucket))
		v := b.Get([]byte(question.Data.Quiz))
		str = string(v)
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
		thirdDb.View(func(thirdTx *bolt.Tx) error {
			// Assume bucket exists and has keys
			b := thirdTx.Bucket([]byte(QuestionBucket))
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				memoryDb.Update(func(tx *bolt.Tx) error {
					b := tx.Bucket([]byte(QuestionBucket))
					err := b.Put(k, v)
					i++
					return err
				})
			}
			log.Println("merged file", f)
			return nil
		})
	}
	log.Println("merged", i, "questions")
}
