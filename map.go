package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Place struct {
	Zip          string
	Abbreviation string
	Latitude     float64
	Longitude    float64
	City         string
	State        string
}

type Locations struct {
	Places []Place
}

func loadLocation(name string) (*Place, error) {
	data, err := ioutil.ReadFile("zips.json")
	if err != nil {
		fmt.Print("Error:", err)
		return nil, err
	}
	var cities Locations
	err = json.Unmarshal(data, &cities)
	if err != nil {
		fmt.Print("Error:", err)
		return nil, err
	}
	for _, elem := range cities.Places {
		if strings.EqualFold(elem.Zip, name) {
			fmt.Printf("%#v\n", elem)
			return &elem, nil
		}
	}
	return nil, nil
}

var templates = template.Must(template.ParseFiles("map.html"))

func mapHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/map/"):]
	p, err := loadLocation(name)
	if err != nil {
		fmt.Print("Error:", err)
		p = &Place{Zip: name}
	}
	err = templates.ExecuteTemplate(w, "map.html", p) 
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func timeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, time.Now().Format("Mon, 02 Jan 2006 15:04:05 MST"))
}

func main() {
	http.HandleFunc("/map/", mapHandler)
	http.HandleFunc("/gettime", timeHandler)
	http.ListenAndServe(":8080", nil)
}
