package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"sync"
)

const (
	html  = `https://www.duitang.com/category/?cat=avatar`
	reImg = `https?://[^"]+?(\.((jpg)|(png)|(jpeg)|(gif)|(bmp)))`
)

var (
	imgUrlChan chan string
	taskChan   chan string
	waitGroup  sync.WaitGroup
)

func handleError(e error, s string) {
	if e != nil {
		log.Fatalln(e, s)
		os.Exit(1)
	}

}

func main() {
	imgUrlChan = make(chan string, 1000)
	taskChan = make(chan string, 26)
	for i := 1; i <= 1; i++ {
		waitGroup.Add(1)
		// mainUrl := html + strconv.Itoa(i) + ".html"
		go getImgUrl(html)
	}
	waitGroup.Add(1)
	go check()
	for i := 1; i <= 5; i++ {
		waitGroup.Add(1)
		go download()
	}
	waitGroup.Wait()
}

func download() {
	for imgUrl := range imgUrlChan {
		index := strings.LastIndex(imgUrl, "/")
		if index == -1 {
			handleError(errors.New("wrong index"), "")
		}
		filename := imgUrl[index+1:]
		filename = "img/" + filename
		r, e := http.Get(imgUrl)
		handleError(e, "get error")
		ct, e := io.ReadAll(r.Body)
		handleError(e, "read error")
		e = os.WriteFile(filename, ct, 0666)
		handleError(e, "download file error")
	}
	waitGroup.Done()
}

func check() {
	var count int
	for url := range taskChan {
		count++
		fmt.Printf("%s has finished the task to capture the imgurl\n", url)
		if count == 26 {
			close(taskChan)
			close(imgUrlChan)
			fmt.Println("all has finished capturing the imgurl")
			break
		}
	}
	waitGroup.Done()
}

func getImgUrl(url string) {
	r, e := http.Get(url)
	handleError(e, "get error!")
	content, e := io.ReadAll(r.Body)
	handleError(e, "read error!")
	contentStr := string(content)
	pattern := regexp.MustCompile(reImg)
	allStr := pattern.FindAllString(contentStr, -1)
	for _, str := range allStr {
		imgUrlChan <- str
		fmt.Println(str)
	}
	taskChan <- url
	waitGroup.Done()
}
