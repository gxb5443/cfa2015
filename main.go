package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"strconv"
)

type Stats struct {
	LandArea          int
	TotalPopulation   int
	TotalHousing      int
	PopulationDensity float32
	HousingDensity    float32
}

type Record struct {
	SummaryLevel        string
	GeographicComponent string
	StateFIPS           string
	PlaceFIPS           string
	CountyFIPS          string
	Tract               string
	Zip                 string
	Block               string
	Name                string
	Latitude            float64
	Longitude           float64
	LandArea            int
	WaterArea           int
	Population          int
	HousingUnits        int
}

func main() {
	csvFile, err := os.Open("./tracts.txt")

	if err != nil {
		fmt.Println(err)
	}

	defer csvFile.Close()

	reader := csv.NewReader(csvFile)

	reader.Comma = '\t'

	reader.FieldsPerRecord = 15

	csvData, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	//var records []*Record
	places := make(map[string]*Stats)
	var maxHDname string
	var maxPDname string
	var maxHD float32
	var maxPD float32
	maxHD = 0.0
	maxPD = 0.0
	var minHDname string
	var minPDname string
	var minHD float32
	var minPD float32
	minHD = 100000.0
	minPD = 100000.0

	for i, row := range csvData {
		if i == 0 {
			continue
		}
		Latitude, _ := strconv.ParseFloat(row[9], 32)
		Longitude, _ := strconv.ParseFloat(row[10], 32)
		LandArea, _ := strconv.Atoi(row[11])
		WaterArea, _ := strconv.Atoi(row[12])
		Population, _ := strconv.Atoi(row[13])
		HousingUnits, _ := strconv.Atoi(row[14])
		r := &Record{
			SummaryLevel:        row[0],
			GeographicComponent: row[1],
			StateFIPS:           row[2],
			PlaceFIPS:           row[3],
			CountyFIPS:          row[4],
			Tract:               row[5],
			Zip:                 row[6],
			Block:               row[7],
			Name:                row[8],
			Latitude:            Latitude,
			Longitude:           Longitude,
			LandArea:            LandArea,
			WaterArea:           WaterArea,
			Population:          Population,
			HousingUnits:        HousingUnits,
		}
		if _, ok := places[r.PlaceFIPS]; !ok {
			places[r.PlaceFIPS] = &Stats{
				LandArea:          r.LandArea,
				TotalPopulation:   r.Population,
				TotalHousing:      r.HousingUnits,
				PopulationDensity: float32(r.Population) / float32(r.LandArea),
				HousingDensity:    float32(r.HousingUnits) / float32(r.LandArea),
			}
			continue
		}
		places[r.PlaceFIPS].TotalHousing += r.HousingUnits
		places[r.PlaceFIPS].TotalPopulation += r.Population
		places[r.PlaceFIPS].PopulationDensity = float32(places[r.PlaceFIPS].TotalPopulation) / float32(places[r.PlaceFIPS].LandArea)
		places[r.PlaceFIPS].HousingDensity = float32(places[r.PlaceFIPS].TotalHousing) / float32(places[r.PlaceFIPS].LandArea)
		if places[r.PlaceFIPS].PopulationDensity > maxPD {
			maxPD = places[r.PlaceFIPS].PopulationDensity
			maxPDname = r.PlaceFIPS
		}
		if places[r.PlaceFIPS].HousingDensity > maxHD {
			maxHD = places[r.PlaceFIPS].HousingDensity
			maxHDname = r.PlaceFIPS
		}
		if places[r.PlaceFIPS].PopulationDensity < maxPD {
			minPD = places[r.PlaceFIPS].PopulationDensity
			minPDname = r.PlaceFIPS
		}
		if places[r.PlaceFIPS].HousingDensity < maxHD {
			minHD = places[r.PlaceFIPS].HousingDensity
			minHDname = r.PlaceFIPS
		}
	}
	fmt.Println("Analysis Complete")
	for k, v := range places {
		fmt.Println("Place: ", k, "----------->")
		fmt.Println("Housing Density: ", v.HousingDensity)
		fmt.Println("Population Density: ", v.PopulationDensity)
	}
	fmt.Println()
	fmt.Println("Max Population Density----->")
	fmt.Println("PlaceFIPS: ", maxPDname, ", Density: ", maxPD)
	fmt.Println()
	fmt.Println("Min Population Density----->")
	fmt.Println("PlaceFIPS: ", minPDname, ", Density: ", minPD)
	fmt.Println("Max housing Density----->")
	fmt.Println("PlaceFIPS: ", maxHDname, ", Density: ", maxHD)
	fmt.Println()
	fmt.Println("Min housing Density----->")
	fmt.Println("PlaceFIPS: ", minHDname, ", Density: ", minHD)
}
