package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

var appID = "qnv4hokm8JK6PbUDAhfv"
var apikey = "ryfRP-zH_4McbF3ercr7YEHv2J7NQYObmSyozW7-1YA"

type output struct {
	Items []struct {
		Title string `json:"title"`
	} `json:"items"`
}

func GetLocation(latitude string, longitude string) string {

	var address string

	url := "https://revgeocode.search.hereapi.com/v1/revgeocode?apiKey=" + apikey + "&at=" + fmt.Sprint(latitude) + "," + fmt.Sprint(longitude)

	res, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var data output
	json.Unmarshal(body, &data)
	// fmt.Println(data)
	for _, add := range data.Items {
		address = add.Title
		fmt.Printf(address)
	}

	return address
}
