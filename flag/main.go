package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"
)

type condition struct {
	Temp string `json:"temp"`
	Text string `json:"text"`
}

type weather struct {
	Query struct {
		Results struct {
			Channel struct {
				Item struct {
					Condition condition `json:"condition"`
				} `json:"item"`
			} `json:"channel"`
		} `json:"results"`
	} `json:"query"`
}

var yahooWeatherUrl = "https://query.yahooapis.com/v1/public/yql?"
var city = flag.String("city", "Bakersfield", "Weather for this city.")
var state = flag.String("state", "CA", "Weather for this state.")

func main() {
	flag.Parse()
	fmt.Println(flag.Args())
	v := url.Values{}
	v.Set("q", fmt.Sprintf("select item.condition from weather.forecast where woeid in (select woeid from geo.places(1) where text=\"%s,%s\")", *city, *state))
	v.Set("format", "json")
	v.Set("env", "store://datatables.org/alltableswithkeys")
	resp, err := http.Get(yahooWeatherUrl + v.Encode())
	if err != nil {
		fmt.Println(err)
	}

	var w weather
	decoder := json.NewDecoder(resp.Body)
	if err := decoder.Decode(&w); err != nil {
		fmt.Println(err)
	}
	c := w.Query.Results.Channel.Item.Condition
	fmt.Printf("The weather in %s, %s is %s with the temperature of %sF.\n", *city, *state, c.Text, c.Temp)
}
