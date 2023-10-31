package main

import (
	"algotrading/asset"
	"algotrading/global"
	"algotrading/logger"
	"fmt"
)

func main() {
	fmt.Println(global.Api_Key)
	fmt.Println(global.Stock_Api)
	s := asset.Stocks{
		Type: asset.Monthly,
		Name: "IBM",
	}
	err := s.Get_Price(10)
	if err != nil {
		logger.Error.Println(err.Error())
		return
	}
	fmt.Println(s.Prices)
}
