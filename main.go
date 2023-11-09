package main

import (
	"algotrading/asset"
	"algotrading/global"
	"algotrading/indicator"
	"algotrading/logger"
	"fmt"
)

func main() {
	fmt.Println(global.Api_Key)
	fmt.Println(global.Stock_Api)

	s := asset.Stocks{
		Type: asset.Daily,
		Name: "IBM",
	}

	err := s.Get_Price(200)

	if err != nil {
		logger.Error.Panic(err.Error())
	}

	fmt.Println(s.Prices)

	sma_indi := indicator.SMA_Indicator{
		Asset_Type: asset.Daily,
		Period:     5,
	}
	sma_indi.Get_Indicator(s)
}
