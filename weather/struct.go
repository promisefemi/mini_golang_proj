package main

import (
	"math"
	"time"
)

type Weather struct {
	ID          int64  `json:"id"`
	Main        string `json:"main"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Current struct {
	Weather []Weather `json:"weather"`
	Temp    float64   `json:"temp"`
	DT      int64     `json:"dt"`
}

func (c Current) ApproxTemp() int {

	return int(math.Round(c.Temp))
}

func (c Current) Date() string {
	current := time.Unix(c.DT, 0)
	return current.Format("03:05PM - Monday, 02 Jan 2006")
}

type Daily struct {
	DT   int64 `json:"dt"`
	Temp struct {
		Day float64 `json:"day"`
	}
	Weather []Weather
}

func (d Daily) ApproxTemp() int {
	return int(math.Round(d.Temp.Day))
}

func (d Daily) Date() string {
	current := time.Unix(d.DT, 0)
	return current.Format("Monday")
}

type Result struct {
	Current     Current
	Daily       []Daily
	CurrentCity City
	Cities      []City
}

type City struct {
	ID      float64
	Name    string
	State   string
	Country string
	Coord   struct {
		Lon float64
		Lat float64
	}
}
