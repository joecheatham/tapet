package main

import (
	"encoding/json"
	"fmt"
	)

type BingAPIResponse struct {
	Images   []struct {
		Startdate     string
		Fullstartdate string
		Enddate       string
		URL           string
		Urlbase       string
		Copyright     string
		Copyrightlink string
		Wp            bool
		Hsh           string
		Drk           int
		Top           int
		Bot           int
		Hs            []struct {
			Desc  string
			Link  string
			Query string
			Locx  int
			Locy  int
		}
		Msg           []interface{}
	}
	Tooltips struct {
				 Loading  string
				 Previous string
				 Next     string
				 Walle    string
				 Walls    string
			 }
}

func getInputImage(body []byte) (string, error) {
	var b = new(BingAPIResponse)
	err := json.Unmarshal(body, &b)
	if (err != nil) {
		fmt.Println("NOPE:", err)
	}
	return "https://bing.com" + b.Images[0].URL, err
}