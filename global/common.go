package global

import (
	"fmt"
	"os"
	"syscall"

	"algotrading/logger"

	"github.com/joho/godotenv"
	"golang.org/x/term"
)

var (
	Algo_Support []string
	Price_Type   []string
)

var (
	Api_Key   string
	Stock_Api string
)

func init() {
	// load .env file and get alphvantage api
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	//从.env中获得API_KEY
	Api_Key = os.Getenv("API_KEY")
	//.env没有api_key，那么我们手动输入
	if Api_Key == "" {
		fmt.Println("==============Input License==============")
		fmt.Println("License:")

		//input key
		key, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			logger.Error.Fatal(err.Error())
		}

		//convert []byte to string
		Api_Key = string(key)
		if Api_Key == "" {
			logger.Error.Fatal("Input error")
		}
	}

	//get alphvantage api
	Stock_Api = os.Getenv("STOCK_API")

	//init Price Type
	Price_Type = append(Price_Type, "Daily")
	Price_Type = append(Price_Type, "Weekly")
	Price_Type = append(Price_Type, "Monthly")

	//init Algo_Support
	Algo_Support = append(Algo_Support, "Stat_Arb")
	Algo_Support = append(Algo_Support, "Mean_Reversion")

}

const (
	Daily = iota + 1
	Weekly
	Monthly

	NoneIndicator
	SMA
	EMA
	MACD
	KDJ

	Stat_Arb
	Mean_Reversion
)
