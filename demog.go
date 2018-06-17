package main

import (
	"fmt"
	"log"
	"os"
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
			state.fips = state.getFIPS()
			// Take that fips id and get the population, household, and median income, stick it in a strut
			fmt.Println(s)
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

func (s USState) getFIPS() int {
	return 0
}

func (s USState) getDemographics(id int) error {
	return nil
}
