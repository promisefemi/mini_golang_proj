package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
)

var api_key string = ""
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
	Title     string
	imageURL  string
	imageData []byte
}

var downloadedImages []*Image

func errorHandler(err interface{}) {
	if err != nil {
		fmt.Println(err)
		fmt.Println("\n")
	}
}

func main() {
	// Initiate Flags
	date := flag.String("date", "", "Date of the image you want to download. Default is today")
	dateRange := flag.String("date-range", "", "Date range for images. e.g Download images from 2020-10-01 to 2020-10-10")
	HD = flag.Bool("hd", true, "Download images in HD. default is true")
	flag.Parse()

	// Check both date and dateRange flags are not empty
	if *dateRange != "" && *date != "" {
		errorHandler("Date and Date Range Flag cannot be present at the same time")
		log.Fatal("")
	}

	// Create a date slice for looping through all dates
	dateSlice := make([]string, 0)
	if *dateRange != "" {
		/*
		   Split date ranges into two. a start date and an end date
		   Substract the end date from the start date to get the number of hours passed,
		   divide hours by 24 to get the number of days,

		   loop through the number of days to add one day to the start date, convert to string and append to the dateSlice array in format YYYY-MM-DD.

		*/
		dates := strings.Split(*dateRange, "to")
		fmt.Println(dates)
		fromDate, _ := time.Parse("2006-01-02", strings.TrimSpace(dates[0]))
		toDate, _ := time.Parse("2006-01-02", strings.TrimSpace(dates[1]))

		if fromDate.After(toDate) {
			log.Fatal("Invalid Date Range: template is [From Date] to [To Date]")
		}
		daysBetween := toDate.Sub(fromDate)
		numberOfDays := int(daysBetween.Hours()) / 24

		for i := 0; i <= int(numberOfDays); i++ {
			currentDay := fromDate.Add(time.Hour * (24 * time.Duration(i)))
			currentDay.String()
			dateSlice = append(dateSlice, currentDay.Format("2006-01-02"))
		}

	} else {
		// Check Date, Convert to string and append to dateSlice
		if *date == "" {
			dateTime := time.Now()
			dateTime.String()
			*date = dateTime.Format("2006-01-02")
		}
		dateSlice = append(dateSlice, *date)
	}

	// Create WaitGroup, number of WaitGroup is the length of the dateSlice

	var wg sync.WaitGroup
	wg.Add(len(dateSlice))

	// Loop dateSlice, makeHTTPCall concurrently and pass the dataChannels to it
	for _, date := range dateSlice {
		if date != "" {
			url := parseURL(date)
			go makeHTTPCalls(url, dataCallChan)
		}
	}

	// Call the Create file Function and First Call Checker concurrently
	go firstCallChecker(&wg)
	go createFile(&wg)

	// Wait for all WaitGroup to complete
	wg.Wait()
}

func firstCallChecker(wg *sync.WaitGroup) {
	// Loop through dataChannel to check for any pending response data
	for data := range dataCallChan {
		// Get data responses process them into nasaResponse Struct

		fmt.Println("Checking File")
		var responseStruct nasaResponse
		err := json.Unmarshal(data, &responseStruct)
		errorHandler(err)

		// Check media type, we are only downloading Images

		if responseStruct.Mediatype == "image" {

			fmt.Println("...Downloading " + responseStruct.Title + " ...")

			// Check URL type
			var imageURL string
			if *HD == true {
				imageURL = responseStruct.HDURL
			} else {
				imageURL = responseStruct.URL
			}
			// Create image struct to save the title and url to be used later during the second image HTTP Call
			image := &Image{Title: responseStruct.Title, imageURL: imageURL}

			downloadedImages = append(downloadedImages, image)

			// make the second HTTP Call to download the images
			go makeHTTPCalls(imageURL, imageCallChan)
		} else {
			// Mark one WaitGroup as done if it is not an image
			wg.Done()
			errorHandler("File type is not an image")
		}

	}

}

func createFile(wg *sync.WaitGroup) {
	for data := range imageCallChan {
		// Loop through the Image channel to find any pending images to write to file

		fmt.Println("Writing File")

		// Compare bytes to find the title of the Image
		title := compareBytes(data)

		// Write the data to file
		err := ioutil.WriteFile(title+".jpg", data, 0644)
		if err != nil {
			errorHandler(err)
		}

		fmt.Println("Success: Image was downloaded successfully")
		// Mark waitGroup as done
		wg.Done()
	}

}

func compareBytes(data []byte) string {
	/* Loop through the downloadImages Slice,
	   compare the bytes to find and return the title of the current Image

	*/
	for _, image := range downloadedImages {
		if correctImage := bytes.Compare(data, image.imageData); correctImage == 0 {
			return image.Title
		}
	}
	return "[Unable to Verify Name]"
}

func makeHTTPCalls(url string, channel chan []byte) {
	// initiate new http Call
	resp, err := http.Get(url)
	defer resp.Body.Close()

	errorHandler(err)

	// passing the response
	body, err := ioutil.ReadAll(resp.Body)
	errorHandler(err)

	// Loop through downloadedImages slice to find any image with the current url,
	// save the []byte to enable comparism white writting to file
	for _, img := range downloadedImages {
		if img.imageURL == url {
			img.imageData = body
		}
		fmt.Println("File Size of Newly Copied")
		fmt.Println(len(img.imageData))

	}

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
