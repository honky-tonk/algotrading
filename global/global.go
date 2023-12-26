package global

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	Api_Key   string
	Stock_Api string
)

func init() {
	fmt.Println("init the env file")
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
	Api_Key = os.Getenv("API_KEY")
	Stock_Api = os.Getenv("STOCK_API")
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
)
