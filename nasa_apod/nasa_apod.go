package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var api_key string = "LlmSFJ90CTQG1coaWaJ3kJmBJOMOQ5UpTQghH06e"
var nasaUrl string = "https://api.nasa.gov/planetary/apod"
var HD *bool
var dataCallChan = make(chan []byte, 10)
var imageCallChan = make(chan []byte, 10)
var done = make(chan bool)

// allImages := ma []
type nasaResponse struct {
	Title     string `json:"title"`
	Mediatype string `json:"media_type"`
	HDURL     string `json:"hdurl"`
	URL       string `json:"url"`
}

type Image struct {
	Title    string
	imageURL string
}

func errorHandler(err interface{}) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("\n")
	}
}

func main() {

	date := flag.String("date", "", "Date of the image you want to download. Default is today")
	dateRange := flag.String("date-range", "", "Date range for images. e.g Download images from 2020-10-01 to 2020-10-10")
	HD = flag.Bool("hd", true, "Download images in HD. default is true")
	flag.Parse()

	if *dateRange != "" && *date != "" {
		errorHandler("Date and Date Range Flag cannot be present at the same time")
		return
	}

	dateSlice := make([]string, 1)
	if *dateRange != "" {

		dates := strings.Split(*dateRange, "to")
		fromDate, _ := time.Parse("2006-01-02", dates[0])
		toDate, _ := time.Parse("2006-01-02", dates[1])

		fmt.Println(fromDate)
		fmt.Println(toDate)

	} else {
		dateSlice = append(dateSlice, *date)
	}

	for _, date := range dateSlice {
		url := parseURL(date)
		go makeHTTPCalls(url, dataCallChan)
	}

	go createFile()
	go firstCallChecker()

	<-done
}

func firstCallChecker() {

	for data := range dataCallChan {

		var responseStruct nasaResponse
		err := json.Unmarshal(data, &responseStruct)
		errorHandler(err)

		if responseStruct.Mediatype == "image" {
			fmt.Println("...Downloading " + responseStruct.Title + " ...")

			var imageURL string
			if *HD == true {
				imageURL = responseStruct.HDURL
			} else {
				imageURL = responseStruct.URL
			}

			go makeHTTPCalls(imageURL, imageCallChan)

		} else {
			errorHandler("File type is not an image")
		}

	}

}

func checkFile() {

}
func createFile() {
	i := 0
	for data := range imageCallChan {

		file, err := os.Create(string(i) + ".jpg")
		if err != nil {
			errorHandler(err)
		}

		_, err = file.Write(data)

		if err != nil {
			errorHandler(err)
		}
		fmt.Println("Success: Image was downloaded successfully")
		if len(imageCallChan) <= 0 {
			done <- true
		}
	}

}

func makeHTTPCalls(url string, channel chan []byte) {

	resp, err := http.Get(url)
	resp.Body.Close()

	errorHandler(err)

	// passing the response
	body, err := ioutil.ReadAll(resp.Body)
	errorHandler(err)
	channel <- body

}

func parseURL(date string) string {
	base, err := url.Parse(nasaUrl)
	errorHandler(err)

	// Adding Query Parameters
	params := url.Values{}
	params.Add("api_key", api_key)
	params.Add("date", date)
	params.Add("hd", strconv.FormatBool(*HD))

	base.RawQuery = params.Encode()

	stringURL := fmt.Sprintf("%v", base.String())

	return stringURL
}
