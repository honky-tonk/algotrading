package main

import (
	"algotrading/asset"
	"algotrading/global"
	"algotrading/indicator"
	"algotrading/logger"
	"fmt"
	//"errors"
)

func test_ema_indicator(s asset.Stocks) {
	ema_indi := indicator.EMA_Indicator{
		Asset_Type: asset.Daily,
		Period:     5,
		Smoothing:  2,
	}

	var err error
	ema_indi.Indicator_Value, err = ema_indi.Get_Indicator(s)

	if err != nil {
		logger.Error.Fatal(err.Error())
	}

}

func main() {
	if global.Api_Key == "" {
		logger.Error.Fatal("Please fill .env's Api_key")
	}
	if global.Stock_Api == "" {
		logger.Error.Fatal("Please fill .env's Stock_Api")
	}
	s := asset.Stocks{
		Type: asset.Daily,
		Name: "IBM",
	}

	err := s.Get_Price(200)

	if err != nil {
		logger.Error.Panic(err.Error())
	}

	fmt.Println(s.Prices)

	//test ema indicator
	test_ema_indicator(s)

}
