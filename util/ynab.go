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

const ynabApiUrl = "https://api.youneedabudget.com/v1"
const ynabApiSingleAccount = "%s/budgets/%s/accounts/%s"

type ynabAccountJsonResponse struct {
	Data struct {
		Account struct {
			Balance float32
		} `json:"account"`
	} `json:"data"`
}

func GetBalanceForAccount(budgetId string, accountId string) float32 {
	if len(budgetId) == 0 || len(accountId) == 0 {
		fmt.Print("Invalid budget or account IDs. Please check your config.")
		os.Exit(1)
	}

	var apiKey = viper.GetString(config.YnabApiKey)

	if len(apiKey) == 0 {
		fmt.Print("Invalid YNAB API Key. Please check your config.")
		os.Exit(1)
	}

	request, _ := http.NewRequest("GET", fmt.Sprintf(ynabApiSingleAccount, ynabApiUrl, budgetId, accountId), nil)
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	response, err := http.DefaultClient.Do(request)

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

	var result ynabAccountJsonResponse
	json.NewDecoder(response.Body).Decode(&result)

	defer response.Body.Close()

	return result.Data.Account.Balance
}
