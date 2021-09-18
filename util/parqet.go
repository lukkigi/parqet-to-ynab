package util

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/lukkigi/parqet-to-ynab/config"
	"github.com/spf13/viper"
)

const apiUrl = "https://api.tresor.one/v1/portfolios/"

type parqetJsonResponse struct {
	Portfolio struct {
		Performance struct {
			PortfolioValue float64
		} `json:"performance"`
	} `json:"portfolio"`
}

func GetPortfolioValue() float64 {
	response, err := http.Get(fmt.Sprintf("%s%s", apiUrl, getPortfolioId()))

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

func getPortfolioId() string {
	portfolioId := viper.GetString(config.ParqetPortfolioId)

	if len(portfolioId) == 0 {
		fmt.Println("Invalid Parqet portfolio ID. Please check your config.")
		os.Exit(1)
	}

	return portfolioId
}
