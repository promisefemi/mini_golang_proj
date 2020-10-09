package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func handleError(name interface{}) {
	if name != nil {

		fmt.Println(name)
		os.Exit(1)
	}
}

func main() {

	noHeader := flag.Bool("noheader", false, "Ignores the first row as title header")
	seperator := flag.String("seperator", "comma", "Default seperator: comma(,) -- only comma(,) and colon(;) are supported")
	flag.Parse()

	if len(flag.Args()) < 1 {
		handleError("No File was passed")

	}

	fileName := flag.Arg(0)

	f, err := os.Open(fileName)
	defer f.Close()
	handleError(err)

	var returnData interface{}

	csvReadFile := csv.NewReader(f)

	if *seperator == "comma" {
		csvReadFile.Comma = ','
	} else if *seperator == "colon" {
		fmt.Println("Colol")
		csvReadFile.Comma = ';'
	} else {
		handleError("Invalid Seperator Passed")
	}

	if *noHeader {
		returnData = parseCSVWithoutHeader(*csvReadFile, *seperator)
	} else {
		returnData = parseCSVWithHeader(*csvReadFile, *seperator)
	}

	// fmt.Println(returnData)

	jsonFile, err := json.Marshal(returnData)
	handleError(err)
	returnFile, err := os.Create(strings.Split(fileName, ".")[0] + ".json")
	handleError(err)

	n, err := returnFile.Write(jsonFile)
	handleError(err)

	fmt.Printf("%d No. of bytes written \n\n", n)
	fmt.Println("SUCCESS")
	// fmt.Println(string(jsonFile))

}

func parseCSVWithoutHeader(csvReadFile csv.Reader, seperator string) [][]string {

	header, err := csvReadFile.Read()
	handleError(err)
	data, err := csvReadFile.ReadAll()
	handleError(err)

	var returnData [][]string

	returnData = append(returnData, header)

	for i := 0; i < len(data); i++ {

		returnData = append(returnData, data[i])
	}

	return returnData

}

func parseCSVWithHeader(csvReadFile csv.Reader, seperator string) []map[string]string {

	header, err := csvReadFile.Read()
	handleError(err)
	data, err := csvReadFile.ReadAll()
	handleError(err)

	var returnData []map[string]string

	for i := 0; i < len(data); i++ {
		arrangeData := make(map[string]string)
		for j := 0; j < len(header); j++ {
			arrangeData[strings.TrimSpace(header[j])] = data[i][j]
		}
		returnData = append(returnData, arrangeData)
	}

	return returnData
}

func checkFileType(fileName string) (bool, error) {

	if filepath.Ext(fileName) == ".csv" {
		return true, nil
	}

	return false, errors.New("Incorrect file type: File is not CSV")

}
