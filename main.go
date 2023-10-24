package main

import (
	"algotrading/asset"
	"algotrading/global"
	"fmt"
)

func main() {
	fmt.Println(global.Api_Key)
	fmt.Println(global.Stock_Api)
	s := asset.Stocks{
		Type: asset.Daily,
		Name: "IBM",
	}
	s.Get_Daily_Price(200)
}
