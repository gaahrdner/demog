package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/urfave/cli"
)

//USState represents a state object
type USState struct {
	name         string
	fips         string
	population   int
	households   int
	medianIncome float64
}

//CensusAPI represents the geographic data from the api
type CensusAPI struct {
	Status       string        `json:"status"`
	ResponseTime int           `json:"responseTime"`
	Message      []interface{} `json:"message"`
	Results      struct {
		State []struct {
			GeographyType string `json:"geographyType"`
			Name          string `json:"name"`
			Fips          string `json:"fips"`
			StateCode     string `json:"stateCode"`
		} `json:"state"`
	} `json:"Results"`
}

//DemographicAPI represents the demographic data from the api
type DemographicAPI struct {
	Status       string        `json:"status"`
	ResponseTime int           `json:"responseTime"`
	Message      []interface{} `json:"message"`
	Results      []struct {
		GeographyID                 string  `json:"geographyId"`
		GeographyName               string  `json:"geographyName"`
		LandArea                    float64 `json:"landArea"`
		Population                  int     `json:"population"`
		Households                  int     `json:"households"`
		RaceWhite                   float64 `json:"raceWhite"`
		RaceBlack                   float64 `json:"raceBlack"`
		RaceHispanic                float64 `json:"raceHispanic"`
		RaceAsian                   float64 `json:"raceAsian"`
		RaceNativeAmerican          float64 `json:"raceNativeAmerican"`
		IncomeBelowPoverty          float64 `json:"incomeBelowPoverty"`
		MedianIncome                float64 `json:"medianIncome"`
		IncomeLessThan25            float64 `json:"incomeLessThan25"`
		IncomeBetween25To50         float64 `json:"incomeBetween25to50"`
		IncomeBetween50To100        float64 `json:"incomeBetween50to100"`
		IncomeBetween100To200       float64 `json:"incomeBetween100to200"`
		IncomeGreater200            float64 `json:"incomeGreater200"`
		EducationHighSchoolGraduate float64 `json:"educationHighSchoolGraduate"`
		EducationBachelorOrGreater  float64 `json:"educationBachelorOrGreater"`
		AgeUnder5                   float64 `json:"ageUnder5"`
		AgeBetween5To19             float64 `json:"ageBetween5to19"`
		AgeBetween20To34            float64 `json:"ageBetween20to34"`
		AgeBetween35To59            float64 `json:"ageBetween35to59"`
		AgeGreaterThan60            float64 `json:"ageGreaterThan60"`
		MyAreaIndicator             bool    `json:"myAreaIndicator"`
	} `json:"Results"`
}

// StateURL is the API endpoint to find geographies specified by a state's name
const StateURL = "https://www.broadbandmap.gov/broadbandmap/census/state/"

// DemographicURL is the API endpoint that returns demographic information
const DemographicURL = "https://www.broadbandmap.gov/broadbandmap/demographic/jun2014/state/ids/"

// Fmt is the format returned from the API
const Fmt = "?format=json"

func main() {

	var output string
	// validOutputs := map[string]bool{"csv": true, "averages": true}

	app := cli.NewApp()
	app.Name = "demog"
	app.Usage = "A CLI for retrieving demographic data for sets of US states"
	app.Author = "Philip Gardner"

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "format, f",
			Usage:       "output format for demographic data [csv,averages]",
			Destination: &output,
		},
	}

	app.Action = func(c *cli.Context) error {

		var data []*USState

		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}

		states := cleanInput(c.Args())

		for _, s := range states {
			// take a state and get the fips id
			state := new(USState)
			state.name = s

			state.fips = getGeoData(s).Results.State[0].Fips

			demo := getDemoData(state.fips)

			state.households = demo.Results[0].Households
			state.population = demo.Results[0].Population
			state.medianIncome = demo.Results[0].MedianIncome

			data = append(data, state)

		}

		if output == "csv" {
			fmt.Println("name,fips,population,households,median_income")
			for _, s := range data {
				fmt.Printf("%v,%v,%d,%d,%f\n", s.name, string(s.fips), s.population, s.households, s.medianIncome)
			}
		} else if output == "averages" {
			fmt.Println(weightedAverage(data))
		} else {
			fmt.Println("The --format flag must be specified and be one of csv or averages")
		}

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}

// cleanInput handles white space issues when inputting states with commas
func cleanInput(a []string) []string {

	// we need to handle white space cases
	arrayToString := strings.Join(a, "")
	array := strings.Split(arrayToString, ",")

	var clean []string
	for _, s := range array {

		s = strings.ToLower(s)
		if strings.Contains(s, "new") {
			clean = append(clean, strings.Replace(s, "new", "new%20", 1))
		} else if strings.Contains(s, "north") {
			clean = append(clean, strings.Replace(s, "north", "north%20", 1))
		} else if strings.Contains(s, "south") {
			clean = append(clean, strings.Replace(s, "south", "south%20", 1))
		} else if strings.Contains(s, "west") {
			clean = append(clean, strings.Replace(s, "west", "west%20", 1))
		} else if strings.Contains(s, "rhode") {
			clean = append(clean, strings.Replace(s, "rhode", "rhode%20", 1))
		} else {
			clean = append(clean, s)
		}
	}

	return clean
}

func getGeoData(s string) *CensusAPI {
	url := StateURL + s + Fmt

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var census CensusAPI

	if err := json.NewDecoder(resp.Body).Decode(&census); err != nil {
		log.Fatal(err)
	}
	return &census
}

func getDemoData(fips string) *DemographicAPI {
	url := fmt.Sprintf(DemographicURL + fips + Fmt)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	var demo DemographicAPI

	if err := json.NewDecoder(resp.Body).Decode(&demo); err != nil {
		log.Fatal(err)
	}

	return &demo
}

func weightedAverage(states []*USState) float64 {
	var sumHouseholds int
	var sumIncome float64
	for _, s := range states {
		sumIncome += s.medianIncome * float64(s.households)
		sumHouseholds += s.households
	}
	return sumIncome / float64(sumHouseholds)
}

// func (s *USState) toArray() []string {
// 	var array []string
// 	array = append(array, s.name)
// 	return array
// }
