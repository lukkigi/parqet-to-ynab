package cmd

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/lukkigi/parqet-to-ynab/util"

	"github.com/briandowns/spinner"
	"github.com/spf13/cobra"
)

var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Syncs your YNAB and Parqet budgets and create an adjusting transaction if needed",
	Long:  `Syncs your YNAB and Parqet budgets and create an adjusting transaction if needed`,
	Run: func(cmd *cobra.Command, args []string) {
		parqetSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		ynabBalanceSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
		ynabTransactionSpinner := spinner.New(spinner.CharSets[11], 100*time.Millisecond)

		parqetSpinner.Start()
		portfolioValue := util.GetPortfolioValue()
		parqetSpinner.FinalMSG = util.StyleText(util.PinkColor, "Portfolio balance", fmt.Sprintf("%.2f", portfolioValue))
		parqetSpinner.Stop()

		ynabBalanceSpinner.Start()
		accountBalance := util.GetBalanceForAccount()
		ynabBalanceSpinner.FinalMSG = util.StyleText(util.PinkColor, "YNAB account balance", fmt.Sprintf("%.2f", accountBalance))
		ynabBalanceSpinner.Stop()

		// YNAB expects the value to be in milliunits format
		difference := int((portfolioValue - accountBalance) * 1000)

		if difference <= 1000 && difference >= -1000 {
			os.Exit(0)
		}

		ynabTransactionSpinner.Start()
		util.AddNewTransaction(strconv.FormatInt(int64(difference), 10))
		ynabTransactionSpinner.FinalMSG = util.StyleText(util.GreenColor, "Transaction created", fmt.Sprintf("%.2f", portfolioValue-accountBalance))
		ynabTransactionSpinner.Stop()
	},
}

func init() {
	rootCmd.AddCommand(syncCmd)
}
