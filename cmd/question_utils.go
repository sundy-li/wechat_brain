//跑这个脚本的时候记得关掉其他服务,因为这个脚本需要读取questions.data
package main

import (
	"archive/zip"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/boltdb/bolt"
	brain "github.com/sundy-li/wechat_brain"

	"github.com/PuerkitoBio/goquery"
)

var (
	source   string
	fs       string
	issueUrl = "https://github.com/sundy-li/wechat_brain/issues/17"
	tmpDir   = "/tmp/"
)

func init() {
	flag.StringVar(&source, "s", "show", "source value -> show | file | issue")
	flag.StringVar(&fs, "fs", "", "merge data files")
	flag.Parse()

	initDb()
}

func main() {
	if source == "file" {
		files := strings.Split(fs, " ")
		if len(files) < 1 {
			log.Println("empty files")
			return
		}
		brain.MergeQuestions(files...)
	} else if source == "issue" {
		doc, _ := goquery.NewDocument(issueUrl)
		doc.Find("div.comment").Each(func(index int, comment *goquery.Selection) {
			comment.Find("td.d-block p a").Each(func(i int, s *goquery.Selection) {
				if strings.Contains(s.Text(), ".zip") {
					href, _ := s.Attr("href")
					if href != "" {
						err := handleZipUrl(href)
						if err != nil {
							log.Println("Error", err.Error())
						}
					}
				}
			})
		})
	} else {
		brain.ShowAllQuestions()
	}
	total := brain.CountQuestions()
	log.Println("total questions =>", total)
}

func handleZipUrl(url string) error {
	println("handling", url)
	var exist bool
	memoryDb.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(QuestionUrlBucket))
		v := b.Get([]byte(url))
		if len(v) != 0 {
			exist = true
		}
		return nil
	})
	if exist {
		log.Println("skip already merged file url", url)
		return nil
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	out, err := os.Create(tmpDir + "questions.zip")
	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, resp.Body)
	_, err = Unzip(tmpDir+"questions.zip", tmpDir)
	if err != nil {
		return err
	}
	//merge data
	brain.MergeQuestions(tmpDir + "questions.data")
	log.Println("merged", url)

	memoryDb.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(QuestionUrlBucket))
		b.Put([]byte(url), []byte("ok"))
		return nil
	})
	return nil
}

// Unzip will un-compress a zip archive,
// moving all files and folders to an output directory
func Unzip(src, dest string) ([]string, error) {
	var filenames []string
	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)
		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			var fdir string
			if lastIndex := strings.LastIndex(fpath, string(os.PathSeparator)); lastIndex > -1 {
				fdir = fpath[:lastIndex]
			}

			err = os.MkdirAll(fdir, os.ModePerm)
			if err != nil {
				log.Fatal(err)
				return filenames, err
			}
			f, err := os.OpenFile(
				fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return filenames, err
			}
			defer f.Close()

			_, err = io.Copy(f, rc)
			if err != nil {
				return filenames, err
			}
		}
	}
	return filenames, nil
}

//questions.data中存入url的cache, 防止重复merge,提高性能

var (
	memoryDb          *bolt.DB
	QuestionUrlBucket = "QuestionUrl"
)

func initDb() {
	var err error
	memoryDb, err = bolt.Open("merge.data", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}
	memoryDb.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte(QuestionUrlBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
}
