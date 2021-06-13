package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"text/template"
)

type WeatherInfo struct {
	Temparature float64 `json:"temp"`
	FeelsLike   float64 `json:"feels_like"`
	MinTemp     float64 `json:"temp_min"`
	MaxTemp     float64 `json:"temp_max"`
	Pressure    float64 `json:"pressure"`
	Humidity    float64 `json:"humidity"`
	SeaLevel    float64 `json:"sea_level"`
	GroundLevel float64 `json:"grnd_level"`
}

type Weather struct {
	Location    string `json:"name"`
	WeatherInfo `json:"main"`
}

// Capitalize first letter and others to lowercase
func CapitalizeInitial(s string) string {
	return strings.ToUpper(string(s[0])) + strings.ToLower(s[1:])
}

func WeatherForecast(w http.ResponseWriter, r *http.Request) {
	tpl := template.Must(template.ParseFiles("index.html"))
	if r.Method == http.MethodPost {
		location := r.FormValue("location")
		location = strings.Trim(location, " ")

		url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s&appid=da8d44e4505f5153cf700b5eeeb1885d&units=metric", CapitalizeInitial(location))

		resp, err := http.Get(url)
		if err != nil {
			http.Error(w, "Invalid location or "+http.StatusText(500), 500)
			log.Fatal("api", err)
		}
		defer resp.Body.Close()

		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal("ReadAll:", err)
		}
		fmt.Println(string(bs))

		fmt.Println("-------------------------------------------")

		var weather Weather
		err = json.Unmarshal(bs, &weather)
		if err != nil {
			log.Fatal("ReadAll:", err)
		}

		fmt.Printf("%+v\n", weather)

		tpl.Execute(w, weather)
		return
	}

	tpl.Execute(w, nil)
}

func main() {
	http.HandleFunc("/", WeatherForecast)

	fmt.Println("visit localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Listenandserve:", err)
	}
}
