package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const apiUrl = "https://api.tresor.one/v1/portfolios/"

type parqetJsonResponse struct {
	Portfolio struct {
		Performance struct {
			PortfolioValue float32
		} `json:"performance"`
	} `json:"portfolio"`
}

func GetPortfolioValue(id string) float32 {
	if len(id) == 0 {
		fmt.Print("Invalid Portfolio ID. Please check your config.")
		os.Exit(1)
	}

	response, err := http.Get(fmt.Sprintf("%s%s", apiUrl, id))

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	if response.StatusCode != http.StatusOK {
		bodyBytes, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(bodyBytes))

		os.Exit(1)
	}

	var result parqetJsonResponse
	json.NewDecoder(response.Body).Decode(&result)

	defer response.Body.Close()

	return result.Portfolio.Performance.PortfolioValue
}
