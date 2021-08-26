package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"strconv"
	"strings"

	//"strconv"
	//"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
)

var title string

func fetchImage(pageUrl string, imageNumber string) {

	resp, err := http.Get(pageUrl)
	if err != nil {
		log.Fatal(err)
	}

	if resp.StatusCode != 200 {
		log.Fatalf("Error: %s", resp.Status)
	}

	defer resp.Body.Close()

	// Initialize HTML Parser

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if resp.StatusCode != 200 {
		log.Fatalf("Error: %s", err)
	}

	var wg sync.WaitGroup

	// project-module grid--main js-grid-main project-module-image-full image-full grid--ready
	// Find HTML Nodes using these selectors (Hopefully they don't change)
	downloadCount := 0

	pageTitle := doc.Find("title")

	title = pageTitle.Text()

	projectModules := doc.Find("#project-modules img")
	projectModules.Each(func(i int, s *goquery.Selection) {
		downloadCount += 1
		if imageNumber == "all" {
			wg.Add(1)
			go processDownload(s, &wg)
		} else {
			for _, number := range strings.Split(imageNumber, ",") {
				if integer, _ := strconv.Atoi(number); integer == i+1 {
					wg.Add(1)
					go processDownload(s, &wg)
				}
			}
		}
	})
	wg.Wait()

	if downloadCount > 0 {
		fmt.Println("Images Downloaded Successfully")
	} else {
		fmt.Println("No Image was Downloaded")
	}
}

func processDownload(s *goquery.Selection, wg *sync.WaitGroup) {
	imgURL, isThere := s.Attr("data-src")

	if !isThere {
		imgURL, isThere = s.Attr("src")
	}
	//fmt.Println(imgURL)

	if isThere {
		fileName := getFileName(imgURL)
		// fmt.Printf("URL PATH = %s , File Name = %s  \n", urlPath, fileName)
		if fileName != "blank.png" {
			fmt.Printf("Downloading: %s\n\n", fileName)
			go downloadImage(fileName, imgURL, wg)
		}
	} else {
		fmt.Println("Could Not find any SRC")
	}
}

func getFileName(imgURL string) string {
	urlPath, _ := url.Parse(imgURL)
	fileName := path.Base(urlPath.Path)
	return fileName
}

func downloadImage(fileName string, imageURL string, wg *sync.WaitGroup) {

	// fmt.Printf("URL of the Image %s \n", imageURL)
	resp, err := http.Get(imageURL)
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		fmt.Printf("%s \n*******------******* \n", resp.Status)
	}

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("%s \n*******------******* \n", err)

	}

	if _, err = os.Stat("assets/" + title + "/"); os.IsNotExist(err) {
		_ = os.Mkdir("assets/"+title, 0755)
		// if err != nil {
		// 	log.Fatal(err)
		// }
	}

	err = ioutil.WriteFile("assets/"+title+"/"+fileName, responseBody, 0777)
	if err != nil {
		fmt.Printf("%s \n*******------******* \n", err)
	}

	fmt.Printf("%s Downloaded \n\n", fileName)

	wg.Done()
}

func main() {
	url := flag.String("url", "", "The Url of the Project")
	imageNumber := flag.String("imagenumber", "all", "Download only selected images (images are counted from 1): Default is 'All'")
	flag.Parse()

	if *url == "" {
		log.Fatal("URL flag is a required field")
	}

	fetchImage(*url, *imageNumber)
}
