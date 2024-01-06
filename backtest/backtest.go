package backtest

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"algotrading/algolib"
	"algotrading/asset"
	"algotrading/global"
)

// for main goroutine pass init signal and data to algo runner goroutine
type AlgoRunner_Init struct {
	Asset_Name      string
	Start_TimePoint time.Time
}

// for init Fetcher, send from algo runner, to fetcher
type Fetcher_Init struct {
	Asset_Name      []string
	Start_TimePoint time.Time
}

// for fetch goroutine fetch data send to algo runner goroutine
type Algo_Message struct {
	Asset_Name string
	SP         asset.Stock_Price
}

// for terminal fetch goroutine and send statistical message to main goroutine, from algo runner goroutine to main goroutine(statistical message)
type Algo_Terminal_And_Statistical struct {
	IsTerminal bool
	Stat       Statistical
}

// algo_runner运行完后的statistical信息
type Statistical struct {
}

func get_asset_names(asset_names []string) error {
	var num int
	//var asset_names []string

	fmt.Println("======Input Number of Assets======")
	fmt.Println("Assets Num: ")
	_, err := fmt.Scan(&num)
	if err != nil {
		return err
	}

	for i := 0; i < num; i++ {
		var asset_name string
		fmt.Println("==========Input Asset=========")
		fmt.Println("Asset", i+1, ": ")
		_, err := fmt.Scan(&asset_name)
		if err != nil {
			return err
		}

		asset_names = append(asset_names, asset_name)
		// assets[i].Name = asset_name
		// assets[i].Period = period
		// assets[i].Type = p_type
		// err = assets[i].Get_Price(db)
		// if err != nil {
		// 	return nil, err
		// }
	}
	return nil

}

// get fetch start timepoint
func get_fetch_start_timepoint(fetch_start_timepoint *string) error {
	fmt.Println("======Input fetch start timepoint======")
	fmt.Println("Start Time: ")
	_, err := fmt.Scan(fetch_start_timepoint)
	if err != nil {
		return err
	}
	return nil
}

// get backtest start time point
func get_backtest_start_timepoint(backtest_start_timepoint *string, fetch_start_timepoint *string) error {
	//.env获取为空然后交互式获取
	if *backtest_start_timepoint == "" {
		fmt.Println("==========Get Backtest Start Time Point========")
		fmt.Println("Start Time Point, like 2006-01-02: ")

		_, err := fmt.Scan(backtest_start_timepoint)
		if err != nil {
			return err
		}
	}

	//backtest_start_point不能小于fetch_start_time_point
	var backtest_start_time time.Time
	var fetch_start_time time.Time

	backtest_start_time, err := time.Parse("2006-01-02", *backtest_start_timepoint)
	if err != nil {
		return err
	}

	fetch_start_time, err = time.Parse("2006-01-02", *fetch_start_timepoint)
	if err != nil {
		return err
	}

	if backtest_start_time.Before(fetch_start_time) {
		return errors.New("error! start time point greater than asset period.")
	}

	return nil
}

// 使用交互的方式get algo
func get_algo(algo *string) error {
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

func get_price_type(price_type_string string) (int, error) {
	//var price_type_string string
	var price_type int
	if price_type_string == "" {
		//.env中没有找到，所以交互式获取
		fmt.Println("Price Type Support Now: \n", "- Daily\n", "- Weekly\n", "- Monthly\n")
		fmt.Println("Inpt Type: ")
		_, err := fmt.Scan(&price_type_string)
		if err != nil {
			return -1, err
		}
	}

	switch {
	case price_type_string == "Daily":
		price_type = 1
		break
	case price_type_string == "Weekly":
		price_type = 2
		break
	case price_type_string == "Monthly":
		price_type = 3
		break
	default:
		return -1, errors.New("Input type of price error")
	}
	return price_type, nil
}

// get price
func get_price(db *sql.DB, assets []asset.Stock, asset_names []string, fetch_timepoint string, price_type int) error {
	start_time, err := time.Parse("2006-01-02", fetch_timepoint)
	if err != nil {
		return err
	}

	for _, v := range asset_names {
		tmp_stock := asset.Stock{
			Name:            v,
			Type:            price_type,
			Start_TimePoint: start_time,
		}
		tmp_stock.Get_Price(db)
		assets = append(assets, tmp_stock)
	}

	return nil
}

// goroutine get each price algo need, and pass to algo goroutine via channel
func fetch_perice() {
	//TODO
}

func exec_algo(algo string, assets []asset.Stock) error {
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
	//用于main goroutine传递初始化给algo runnter goroutine
	AlgoRunner_Init_Chan := make(chan AlgoRunner_Init)
	//用于algo runner发送初始化信息给Fetcher goroutine
	Fetcher_Init_Chan := make(chan Fetcher_Init)
	//用于Fetcher发送message给algo runner
	Algo_message_Chan := make(chan []Algo_Message)
	//用于algo runner goroutine发送terminal信号给Fetcher goroutine，和Statistical给main
	Algo_Ter_Stat_Chan := make(chan Algo_Terminal_And_Statistical)

	//从.env中获取algo
	algo := os.Getenv("ALGO")
	if algo == "" {
		//如果.env中没有找到ALGO那么手动获取
		err := get_algo(&algo)
		if err != nil {
			return err
		}
	}

	//从.env中获取asset的所有name
	asset_names_row := os.Getenv("ASSETS")
	var asset_names []string
	if asset_names_row == "" {
		//如果.env没有找到就交互获取
		err := get_asset_names(asset_names)
		if err != nil {
			return err
		}
	} else {
		for _, v := range strings.Split(asset_names_row, ",") {
			asset_names = append(asset_names, v)
		}
	}

	//从.env中获取asset的type
	var price_type_string string
	var price_type int
	price_type_string = os.Getenv("PRICE_TYPE")
	//无论在.env没有找到都要进这个函数转换成int类型
	price_type, err := get_price_type(price_type_string)
	if err != nil {
		return err
	}

	//从.env中获取fetch start time point
	var fetch_start_timepoint string
	fetch_start_timepoint = os.Getenv("FETCH_START_TIMEPOINT")
	if fetch_start_timepoint == "" {
		//如果.env没有找到就交互获取
		err := get_fetch_start_timepoint(&fetch_start_timepoint)
		if err != nil {
			return err
		}
	}

	//从.env中获取backtest_start_time_point
	var backtest_start_timepoint string
	backtest_start_timepoint = os.Getenv("BACKTEST_START_TIMEPOINT")
	//无论是否从.env获取到了BACKTEST_START_TIMEPOINT，都要进入下面函数进行判断
	err = get_backtest_start_timepoint(&backtest_start_timepoint, &fetch_start_timepoint)
	if err != nil {
		return err
	}

	//get price
	var assets []asset.Stock
	err = get_price(db, assets, asset_names, fetch_start_timepoint, price_type)
	if err != nil {
		return err
	}

	//exec algo
	//先说一下这里的思想，由一个goroutine去哪一天的数据然后传入channel被阻塞，
	//algo也是由一个goroutine驱动，algo goroutine从channel拿到数据进行计算，
	//在algo goroutine拿到数据的时候fetch_price goroutine解除阻塞继续执行
	exec_algo(algo, assets, start_timepoint, AlgoRunner_Init_Chan)

	//statistical of backtest

	return nil

}
