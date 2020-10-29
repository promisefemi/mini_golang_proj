package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path"
	"sync"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func fetchImage(pageUrl string) {

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
		imgUrl, isThere := s.Attr("src")
		if isThere {
			urlPath, _ := url.Parse(imgUrl)

			fileName := path.Base(urlPath.Path)

			if fileName != "blank.png" {
				wg.Add(1)
				fmt.Printf("Downloading: %s\n", fileName)
				go downloadImage(fileName, imgUrl, &wg)
			}
		} else {
			fmt.Println("Could Not find any SRC")
		}

	})
	wg.Wait()
	fmt.Println("Images Downloaded Successfully")

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
	
	if _, err = os.Stat("assets/"); os.IsNotExist(err){
		err = os.Mkdir("assets",0755)
		if err != nil{
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
	flag.Parse()

	if *url == "" {
		log.Fatal("URL flag is a required field")
	}

	fetchImage(*url)
}
