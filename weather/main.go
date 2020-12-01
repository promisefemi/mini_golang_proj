package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"text/template"
)

var appid = "" //Replace the api key
var baseURL = "http://api.openweathermap.org/data/2.5/onecall"
var tpl = template.New("index.html")
var currentCity = City{}
var cities = []City{}

func indexHandler(w http.ResponseWriter, r *http.Request) {

	cityid := float64(2332453)

	u, _ := url.Parse(r.URL.String())

	qParams := u.Query()

	if id := qParams.Get("cityid"); id != "" {
		conv, _ := strconv.Atoi(id)
		cityid = float64(conv)
	}

	parseCities(cityid)

	if currentCity == (City{}) {

		data := make(map[string]interface{})
		data["Cities"] = cities

		err := tpl.Execute(w, data)
		if err != nil {
			log.Println(err)
		}
	}

	params := make(map[string]string)

	params["lat"] = strconv.Itoa(int(currentCity.Coord.Lat))
	params["lon"] = strconv.Itoa(int(currentCity.Coord.Lon))
	params["exclude"] = "minutely,hourly,alerts"

	url := parseURL(params)

	response := makeHTTPCalls(w, url)

	responseData := &Result{}
	err := json.Unmarshal(response, &responseData)
	reportError(w, err)
	responseData.CurrentCity = currentCity
	responseData.Cities = cities

	// empJSON, err := json.MarshalIndent(responseData, "", "  ")
	// if err != nil {
	// 	log.Fatalf(err.Error())
	// }
	// fmt.Printf("MarshalIndent funnction output %s\n", string(empJSON))

	tpl = template.Must(tpl.Funcs(template.FuncMap{
		"toString": toString,
	}).ParseFiles("index.html", "modal.html"))

	err = tpl.Execute(w, responseData)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("./assets/"))
	mux.Handle("/assets/", http.StripPrefix("/assets/", fs))
	mux.HandleFunc("/", indexHandler)
	fmt.Println(" ----- Ready ----- ")

	log.Fatal(http.ListenAndServe(":8080", mux))

}

func makeHTTPCalls(w http.ResponseWriter, url string) []byte {

	resp, err := http.Get(url)
	reportError(w, err)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	reportError(w, err)

	return body
}

func reportError(w http.ResponseWriter, err error) {
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Fatal(err)
	}
}

func parseCities(cityID float64) {

	data, err := ioutil.ReadFile("list.json")
	if err != nil {
		log.Fatal("Unable to Read City file")
	}
	// fmt.Printf("%s\n", data)
	err = json.Unmarshal(data, &cities)
	if err != nil {
		log.Fatal(err.Error() + " Unable to Parse Cities to Struct")
	}

	for _, city := range cities {
		if city.ID == cityID {
			currentCity = city
			break
		}
	}

}

func parseURL(query map[string]string) string {
	baseURL, _ := url.Parse(baseURL)

	params := url.Values{}
	params.Add("appid", appid)
	params.Add("units", "metric")

	for k, v := range query {
		params.Add(k, v)
	}

	baseURL.RawQuery = params.Encode()

	return fmt.Sprintf("%v", baseURL.String())
}

func toString(cities []City) string {
	// fmt.Println("ajsdkfasdfsk")
	data, err := json.Marshal(cities)
	if err != nil {
		return "[]"
	}

	return string(data)

}
