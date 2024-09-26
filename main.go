package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	url2 "net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
)

var (
	urlFlag string
)

func init() {
	flag.StringVar(&urlFlag, "u", "", "Book link from books.com.tw")
}

func main() {
	flag.Parse()
	//fmt.Println("urlFlag:", urlFlag)
	if len(urlFlag) == 0 {
		log.Fatal("Include the -u flag to specify the source for fetching book cover.")
	}

	req, err := http.NewRequest("GET", urlFlag, nil)
	if err != nil {
		log.Fatal(err)
	}

	// Do GET request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Find cover image
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	title := doc.Find("h1").First().Text()
	fmt.Println("Book Title:", title)

	rawImageUrl, exists := doc.Find("img.cover").First().Attr("src")
	if exists == false {
		log.Fatal(err)
	}
	// fmt.Println("rawImageUrl:", rawImageUrl)

	url, err := url2.Parse(rawImageUrl)
	if err != nil {
		log.Fatal(err)
	}

	srcImgUrl := url.Query()["i"][0]
	fmt.Println("Book Cover:", srcImgUrl)

	saveImageFromUrl(title, srcImgUrl)
}

func saveImageFromUrl(title string, url string) {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(title + ".jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Save book cover success!")
}
