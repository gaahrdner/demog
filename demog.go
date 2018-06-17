package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/urfave/cli"
)

//USState represents a state object
type USState struct {
	name         string
	fips         int
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

//DemographicAPI
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
const DemographicURL = "https://www.broadbandmap.gov/broadbandmap/demographic/jun2014/"

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

		if len(c.Args()) == 0 {
			cli.ShowAppHelp(c)
			os.Exit(0)
		}

		states := cleanInput(c.Args())

		for _, s := range states {
			// take a state and get the fips id
			state := new(USState)
			state.name = s

			fips, err := strconv.Atoi(getGeoData(s).Results.State[0].Fips)
			if err != nil {
				log.Fatal(err)
			}
			state.fips = fips
			fmt.Println(state)
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

	var dirty, clean []string
	if len(a) == 1 {
		dirty = strings.Split(a[0], ",")
	} else {
		dirty = a
	}
	for _, state := range dirty {
		clean = append(clean, strings.Trim(state, ","))
	}
	return clean
}

func getGeoData(s string) *CensusAPI {
	url := fmt.Sprintf(StateURL + s + Fmt)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	defer resp.Body.Close()

	var census CensusAPI

	if err := json.NewDecoder(resp.Body).Decode(&census); err != nil {
		log.Println(err)
	}

	return &census
}
