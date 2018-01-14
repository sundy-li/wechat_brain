package main

import (
	"archive/zip"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
				if strings.Contains(s.Text(), "questions.zip") {
					href, _ := s.Attr("href")
					if href != "" {
						handleZipUrl(href)
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
