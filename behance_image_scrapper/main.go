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
	"sync"

	"github.com/PuerkitoBio/goquery"
)

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

	// Find HTML Nodes using these selectors (Hopefully they don't change)
	doc.Find("#project-modules .project-module-image-inner-wrap img").Each(func(i int, s *goquery.Selection) {

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
	fmt.Println("Images Downloaded Successfully")

}

func processDownload(s *goquery.Selection, wg *sync.WaitGroup) {
	imgUrl, isThere := s.Attr("src")
	if isThere {
		urlPath, _ := url.Parse(imgUrl)

		fileName := path.Base(urlPath.Path)
		if fileName != "blank.png" {
			fmt.Printf("Downloading: %s\n\n", fileName)
			go downloadImage(fileName, imgUrl, wg)
		}
	} else {
		fmt.Println("Could Not find any SRC")
	}

}

func downloadImage(fileName string, imageURL string, wg *sync.WaitGroup) {

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

	if _, err = os.Stat("assets/"); os.IsNotExist(err) {
		err = os.Mkdir("assets", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = ioutil.WriteFile("assets/"+fileName, responseBody, 0777)
	if err != nil {
		fmt.Printf("%s \n*******------******* \n", err)
	}

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
