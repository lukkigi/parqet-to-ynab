package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lukkigi/parqet-to-ynab/config"
	"github.com/spf13/viper"
)

const ynabApiUrl = "https://api.youneedabudget.com/v1"
const ynabApiSingleAccount = "%s/budgets/%s/accounts/%s"
const ynabApiCreateTransaction = "%s/budgets/%s/transactions"

type ynabAccountJsonResponse struct {
	Data struct {
		Account struct {
			Balance float64
		} `json:"account"`
	} `json:"data"`
}

type ynabNewTransactionBody struct {
	Transaction ynabTransaction `json:"transaction"`
}

type ynabTransaction struct {
	AccountId string `json:"account_id"`
	Date      string `json:"date"`
	Amount    string `json:"amount"`
	PayeeName string `json:"payee_name"`
	Approved  bool   `json:"approved"`
	Cleared   string `json:"cleared"`
}

func GetBalanceForAccount() float64 {
	apiKey := getApiKey()
	budgetId := getBudgetId()
	accountId := getAccountId()

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

	return result.Data.Account.Balance / 1000
}

func AddNewTransaction(amount string) string {
	apiKey := getApiKey()
	budgetId := getBudgetId()
	accountId := getAccountId()

	date := time.Now().Local()

	transaction := ynabNewTransactionBody{
		ynabTransaction{
			AccountId: accountId,
			Date:      fmt.Sprintf("%d-%d-%d", date.Year(), date.Month(), date.Day()),
			Amount:    amount,
			PayeeName: "Parqet2YNAB",
			Approved:  true,
			Cleared:   "cleared",
		},
	}

	jsonValue, _ := json.Marshal(transaction)

	request, _ := http.NewRequest("POST", fmt.Sprintf(ynabApiCreateTransaction, ynabApiUrl, budgetId), bytes.NewBuffer(jsonValue))
	request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", apiKey))
	request.Header.Set("Content-Type", "application/json")
	response, err := http.DefaultClient.Do(request)

	if err != nil {
		fmt.Print(err.Error())
		os.Exit(1)
	}

	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}

	return string(bodyBytes)
}

func getApiKey() string {
	apiKey := viper.GetString(config.YnabApiKey)

	if len(apiKey) == 0 {
		fmt.Println("Invalid YNAB API Key. Please check your config.")
		os.Exit(1)
	}

	return apiKey
}

func getAccountId() string {
	accountId := viper.GetString(config.YnabInvestingAccountId)

	if len(accountId) == 0 {
		fmt.Println("Invalid YNAB account ID. Please check your config.")
		os.Exit(1)
	}

	return accountId
}

func getBudgetId() string {
	budgetId := viper.GetString(config.YnabBudgetId)

	if len(budgetId) == 0 {
		fmt.Println("Invalid YNAB budget ID. Please check your config.")
		os.Exit(1)
	}

	return budgetId
}
