package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Person struct {
	Name                    string  `json:"name"`
	Sex                     string  `json:"sex"`
	Age                     int64   `json:"age"`
	PassengerClass          int64   `json:"passengerClass"`
	SiblingsOrSpousesAboard int64   `json:"siblingsOrSpousesAboard"`
	ParentsOrChildrenAboard int64   `json:"parentsOrChildrenAboard"`
	Fare                    float64 `json:"fare"`
	Survived                bool    `json:"survived"`
}

func main() {
	path := os.Args[1]
	host := os.Args[2]
	port := os.Args[3]

	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Fatal("Unable to open file "+path, err)
	}

	csvReader := csv.NewReader(f)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatal("Unable to parse file as CSV for "+path, err)
	}
	for i := range records {
		//skip headers
		if i == 0 {
			continue
		}

		r := records[i]
		if len(r) != 8 {
			log.Printf("Row doesn't contain the right number of fields. Skipping.")
			continue
		}
		survived, err := strconv.ParseBool(r[0])
		if err != nil {
			log.Printf("Couldn't parse 'Survived' field for value %s\n", r[0])
			continue
		}
		passengerClass, err := strconv.ParseInt(r[1], 10, 64)
		if err != nil {
			log.Printf("Couldn't parse 'PassengerClass' field for value %s\n", r[1])
			continue
		}
		name := r[2]
		sex := r[3]
		//age in csv file is a float but the API should only accept ints
		ageFloat, err := strconv.ParseFloat(r[4], 64)
		age := int64(ageFloat)
		if err != nil {
			log.Printf("Couldn't parse 'Age' field for value %s\n", r[4])
			continue
		}
		siblingsOrSpousesAboard, err := strconv.ParseInt(r[5], 10, 64)
		if err != nil {
			log.Printf("Couldn't parse 'SiblingsOrSpousesAboard' field for value %s\n", r[5])
			continue
		}
		parentsOrChildrenAboard, err := strconv.ParseInt(r[6], 10, 64)
		if err != nil {
			log.Printf("Couldn't parse 'ParentsOrChildrenAboard' field for value %s\n", r[6])
			continue
		}
		fare, err := strconv.ParseFloat(r[7], 64)
		if err != nil {
			log.Printf("Couldn't parse 'Fare' field for value %s\n", r[7])
			continue
		}
		p := Person{
			Name:                    name,
			Sex:                     sex,
			Age:                     age,
			PassengerClass:          passengerClass,
			SiblingsOrSpousesAboard: siblingsOrSpousesAboard,
			ParentsOrChildrenAboard: parentsOrChildrenAboard,
			Fare:                    fare,
			Survived:                survived,
		}

		uri := fmt.Sprintf("http://%s:%s/people", host, port)
		body, err := json.Marshal(p)
		if err != nil {
			log.Fatal("Unable to marshal JSON asrequest body")
		}
		res, err := http.Post(uri, "application/json", bytes.NewBuffer(body))
		if err != nil {
			log.Printf("Error while sending http request: %s\n", err)
			continue
		}
		if res.Status == "200 OK" {
			log.Printf("Created person %s\n", p.Name)
		} else {
			log.Printf("Error creating person %s: Server returned %s", p.Name, res.Status)
		}
	}
}
