package backtest

import (
	"database/sql"
	"errors"
	"fmt"

	"algotrading/algolib"
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

// get backtest start time point
func select_backtest_start_timepoint(asset_period int) (int, error) {
	fmt.Println("==========Get Backtest Start Time Point========")
	fmt.Println("Start Time Point: ")

	var start_period int
	_, err := fmt.Scan(&start_period)
	if err != nil {
		return 0, err
	}

	if start_period >= asset_period {
		return 0, errors.New("error! start time point greater than asset period.")
	}

	return start_period, nil
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

// goroutine get each price algo need, and pass to algo goroutine via channel
func fetch_perice() {
	//TODO
}

func exec_algo(algo string, assets []asset.Stocks) error {
	switch {
	case algo == "Stat_Arb":
		params := algolib.Params{}
		params.IsBackTest = true
		params.S = assets
		algolib.Call_Algo(params, algolib.Stat_Arb)

	case algo == "Mean_Reversion":

	default:
		return errors.New("Algo Input Error: Algo Not Found")
	}
	return nil
}

func statistical_of_backtest() {

}

func Backtest_Main(db *sql.DB) error {
	//在回测的时候每一次都由fetch goroutine发送asset.Price给
	message_ := make(chan asset.Backtest_Messages)
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

	//select backtest start time point
	start_timepoint, err := select_backtest_start_timepoint(assets[0].Period)
	if err != nil {
		return err
	}

	//exec algo
	//先说一下这里的思想，由一个goroutine去哪一天的数据然后传入channel被阻塞，
	//algo也是由一个goroutine驱动，algo goroutine从channel拿到数据进行计算，
	//在algo goroutine拿到数据的时候fetch_price goroutine解除阻塞继续执行
	exec_algo(algo, assets)

	//statistical of backtest

	return nil

}
