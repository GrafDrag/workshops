package main

import (
	"encoding/json"
	"log"
	"os"
	"sort"
)

type Person struct {
	name string
	age  int
}

type ByAge []Person

func (a ByAge) Len() int {
	return len(a)
}

func (a ByAge) Less(i, j int) bool {
	return a[i].age < a[j].age
}

func (a ByAge) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type Place struct {
	city    string
	country string
}

type ByCityName []Place

func (a ByCityName) Len() int {
	return len(a)
}

func (a ByCityName) Less(i, j int) bool {
	return len(a[i].city) < len(a[j].city)
}

func (a ByCityName) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

type JsonData struct {
	Things []struct {
		Name    string `json:"name"`
		Age     int    `json:"age"`
		City    string `json:"city"`
		Country string `json:"country"`
	} `json:"things"`
}

type HumanDecoder interface {
	Decode(data []byte) ([]Person, []Place)
	Sort(dataToSort interface{})
	Print(interface{})
}

type Logger interface {
	Println(v ...interface{})
	Fatalf(format string, v ...interface{})
}

type Service struct {
	log Logger
}

func (s Service) Decode(data []byte) ([]Person, []Place) {
	persons := []Person{}
	places := []Place{}
	var res JsonData

	if err := json.Unmarshal(data, &res); err != nil {
		s.log.Fatalf("problem parse json: %v", err)
	}

	for _, data := range res.Things {
		switch {
		case data.Age != 0 && data.Name != "":
			persons = append(persons, Person{
				name: data.Name,
				age:  data.Age,
			})
		case data.City != "" && data.Country != "":
			places = append(places, Place{
				city:    data.City,
				country: data.Country,
			})
		}
	}

	return persons, places
}

func (s Service) Sort(dataToSort interface{}) {
	if data, ok := dataToSort.([]Person); ok {
		sort.Sort(ByAge(data))
	}
	if data, ok := dataToSort.([]Place); ok {
		sort.Sort(ByCityName(data))
	}
}

func (s Service) Print(data interface{}) {
	s.log.Println(data)
}

func main() {
	// logger to Inject
	logger := log.New(os.Stdout, "INFO: ", 0)

	service := Service{
		log: logger,
	}

	persons, places := service.Decode(jsonStr)

	service.Sort(persons)
	service.Sort(places)

	service.Print(persons)
	service.Print(places)
}

var jsonStr = []byte(`
{
    "things": [
        {
            "name": "Alice",
            "age": 37
        },
        {
            "city": "Ipoh",
            "country": "Malaysia"
        },
        {
            "name": "Bob",
            "age": 36
        },
        {
            "city": "Northampton",
            "country": "England"
        },
 		{
            "name": "Albert",
            "age": 3
        },
		{
            "city": "Dnipro",
            "country": "Ukraine"
        },
		{
            "name": "Roman",
            "age": 32
        },
		{
            "city": "New York City",
            "country": "US"
        }
    ]
}`)
