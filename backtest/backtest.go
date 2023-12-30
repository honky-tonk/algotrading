package backtest

import (
	"database/sql"
	"errors"
	"fmt"

	"algotrading/asset"
	"algotrading/global"
)

func select_asset(db *sql.DB) ([]asset.Stocks, error) {
	var num int
	var period int
	var p_type int
	var price_type string

	fmt.Println("======Input Number of Assets======")
	fmt.Println("Assets Num: ")
	_, err := fmt.Scan(&num)
	if err != nil {
		return nil, err
	}

	fmt.Println("======Input Period of Assets======")
	fmt.Println("Period: ")
	_, err = fmt.Scan(&period)
	if err != nil {
		return nil, err
	}

	fmt.Println("=======Input Type of Assets(number)=======")
	for _, v := range global.Price_Type {
		fmt.Println("- ", v)
	}
	fmt.Println("Type: ")
	_, err = fmt.Scan(&price_type)
	if err != nil {
		return nil, err
	}
	switch {
	case price_type == "Daily":
		p_type = 1
		break
	case price_type == "Weekly":
		p_type = 2
		break
	case price_type == "Monthly":
		p_type = 3
		break
	default:
		return nil, errors.New("Input type of price error")
	}

	assets := make([]asset.Stocks, num)
	for i := 0; i < num; i++ {
		var asset_name string
		fmt.Println("==========Input Asset=========")
		fmt.Println("Asset", i+1, ": ")
		_, err := fmt.Scan(&asset_name)
		if err != nil {
			return nil, err
		}

		assets[i].Name = asset_name
		assets[i].Period = period
		assets[i].Type = p_type
		err = assets[i].Get_Price(db)
		if err != nil {
			return nil, err
		}
	}
	return assets, nil

}

func select_algo(algo *string) error {
	//print algo support now
	fmt.Println("==========Algo Support Now========")
	for _, v := range global.Algo_Support {
		fmt.Println("- ", v)
	}

	//select algo support now
	fmt.Println("=============Select Algo===========")
	fmt.Println("Algo: ")
	_, err := fmt.Scan(algo)
	if err != nil {
		return err
	}
	return nil
}

func exec_algo(algo string, assets []asset.Stocks) error {
	switch {
	case algo == "Stat_Arb":

	case algo == "Mean_Reversion":

	default:
		return errors.New("Algo Input Error: Algo Not Found")
	}
	return nil
}

func statistical_of_backtest() {

}

func Backtest_Main(db *sql.DB) error {
	//select algo
	var algo string
	err := select_algo(&algo)
	if err != nil {
		return err
	}

	//select assets
	assets, err := select_asset(db)
	if err != nil {
		return err
	}

	//exec algo
	exec_algo(algo, assets)

	//statistical of backtest

	return nil

}
