package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"strings"
)

type Place struct {
	Name string
	Lat  float64
	Lng  float64
}

type Locations struct {
	Places []Place
}

func loadPlace(name string) (*Place, error) {
	filename := name + ".json"
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Print("Error:", err)
		return nil, err
	}
	var place Place
	err = json.Unmarshal(data, &place)
	if err != nil {
		fmt.Print("Error:", err)
		return nil, err
	}
	fmt.Printf("%#v\n", place)
	return &place, nil
}

func loadLocation(name string) (*Place, error) {
	data, err := ioutil.ReadFile("cities.json")
	if err != nil {
		return nil, err
	}
	var cities Locations
	err = json.Unmarshal(data, &cities)
	if err != nil {
		return nil, err
	}
	for _, elem := range cities.Places {
		if strings.EqualFold(elem.Name, name) {
			fmt.Printf("%#v\n", elem)
			return &elem, nil
		}
	}
	return nil, nil
}

func mapHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/map/"):]
	p, err := loadLocation(name)
	if err != nil {
		fmt.Print("Error:", err)
		p = &Place{Name: name}
	}
	t, _ := template.ParseFiles("map.html")
	t.Execute(w, p)
}

func main() {
	http.HandleFunc("/map/", mapHandler)
	http.ListenAndServe(":8080", nil)
}
