package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"os"
)

type AllApiData struct {
	Images  []string
	Weather *WeatherData
}

type WeatherData struct {
	Temp string
	City string
	RainShine string
}

var dispdata AllApiData
var celsiusNum string
var imagesArray []string
var rainOrShine string
var randIndex int

var cityLibrary = []string{"Tokyo", "Paris", "Singapore", "Sendai", "London", "Shanghai", "Beijing", "Seoul", "Mumbai", "Washington", "Bangkok", "Hanoi", "Toronto", "Atlanta", "Rome", "Milan", "Edinburgh", "Vienna", "Prague", "Stockholm", "Vancouver", "Barcelona", "Sydney", "Istanbul", "Hokkaido", "Santiago", "Valencia", "Peru", "Moscow", "Florence", "Berlin", "Auckland", "Kyoto"}

func ImageDisplay() {
	reqUrl := fmt.Sprintf("https://api.flickr.com/services/rest/?method=flickr.photos.search&api_key=%s&tags=%s&extras=url_m&format=json&nojsoncallback=1&min_taken_date=1388534400&sort=relevance", os.Getenv("FLICKR_APIKEY"), cityLibrary[randIndex])

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
		return
	}
	defer resp.Body.Close()

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}

	type FlickrResponse struct {
		Photos struct {
			Photo []struct {
				Id, Secret, Server string
				Farm               int
			}
		}
	}

	var f FlickrResponse
	errr := json.Unmarshal(body, &f)
	if errr != nil {
		log.Fatal(errr)
	}
	imagesArray = []string{}
	for i := 0; i < 27; i++ {
		v := rand.Intn(100)
		respUrl := "https://farm" + strconv.Itoa(f.Photos.Photo[v].Farm) + ".staticflickr.com/" + f.Photos.Photo[v].Server + "/" + f.Photos.Photo[v].Id + "_" + f.Photos.Photo[v].Secret + "_q.jpg"
		imagesArray = append(imagesArray, respUrl)
	}
}

func WeatherDisplay() {
	reqUrl := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?q=%s", cityLibrary[randIndex])

	client := &http.Client{}
	req, err := http.NewRequest("GET", reqUrl, nil)
	if err != nil {
		log.Fatal("NewRequest: ", err)
		return
	}
	resp, requestErr := client.Do(req)
	if requestErr != nil {
		log.Fatal("Do: ", requestErr)
		return
	}
	defer resp.Body.Close()

	body, dataReadErr := ioutil.ReadAll(resp.Body)
	if dataReadErr != nil {
		log.Fatal("ReadAll: ", dataReadErr)
		return
	}

	type WeatherResponse struct {
		Main struct {
			Temp float64
		}
		Weather []struct {
			Icon string
		}
	}
	var f WeatherResponse
	errr := json.Unmarshal(body, &f)
	if errr != nil {
		log.Fatal(errr)
	}

	celsiusNum = fmt.Sprintf("%.1f", f.Main.Temp-273.15)
	rainOrShine = fmt.Sprintf("http://openweathermap.org/img/w/%s.png", f.Weather[0].Icon)
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	randIndex = rand.Intn(len(cityLibrary))
	ImageDisplay()
	WeatherDisplay()
	dispdata = AllApiData{Images: imagesArray, Weather: &WeatherData{Temp: celsiusNum, City: cityLibrary[randIndex], RainShine: rainOrShine}}
	renderTemplate(w, "home", dispdata)
	}

func renderTemplate(w http.ResponseWriter, tmpl string, structdata AllApiData) {
	t := template.Must(template.New("image").ParseFiles("layout/home.html"))
	tErr := t.ExecuteTemplate(w, tmpl, structdata)
	if tErr != nil {
		http.Error(w, tErr.Error(), http.StatusInternalServerError)
	}
}

func main() {
	http.HandleFunc("/", homeHandler)

	http.Handle("/layout/", http.StripPrefix("/layout/", http.FileServer(http.Dir("layout"))))

	http.ListenAndServe(os.Getenv("PORT"), nil)
}
