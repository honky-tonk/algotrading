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
	Algo                     string
	Assets                   []asset.Stock
	Backtest_Start_TimePoint time.Time
}

// for init Fetcher, send from algo runner, to fetcher
type Fetcher_Init struct {
	Asset_Names     []string
	Start_TimePoint time.Time
}

// for fetch goroutine fetch data send to algo runner goroutine
type Algo_Message struct {
	Asset_Name string
	P          asset.Price
}

// for terminal fetch goroutine and send statistical message to main goroutine, from algo runner goroutine to main goroutine(statistical message)
type Algo_Terminal_And_Statistical struct {
	IsTerminal bool
	Stat       Statistical
}

type Err_Message struct {
	Gorotuine_Type string
	Err            error
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
func get_fetch_start_timepoint(fetch_start_timepoint_raw string, fetch_start_timepoint *time.Time) error {
	if fetch_start_timepoint_raw == "" {
		fmt.Println("======Input fetch start timepoint======")
		fmt.Println("Start Time: ")
		_, err := fmt.Scan(&fetch_start_timepoint_raw)
		if err != nil {
			return err
		}
		return nil
	}

	var time_point time.Time
	time_point, err := time.Parse("2006-01-02", fetch_start_timepoint_raw)
	if err != nil {
		return err
	}
	*fetch_start_timepoint = time_point
	return nil

}

// get backtest start time point
func get_backtest_start_timepoint(backtest_start_timepoint_raw string, backtest_start_timepoint *time.Time, fetch_start_timepoint time.Time) error {
	//.env获取为空然后交互式获取
	if backtest_start_timepoint_raw == "" {
		fmt.Println("==========Get Backtest Start Time Point========")
		fmt.Println("Start Time Point, like 2006-01-02: ")
		_, err := fmt.Scan(&backtest_start_timepoint_raw)
		if err != nil {
			return err
		}
	}

	var timepoint time.Time
	//parse get time string
	timepoint, err := time.Parse("2006-01-02", backtest_start_timepoint_raw)
	*backtest_start_timepoint = timepoint
	if err != nil {
		return err
	}

	//backtest_start_point不能小于fetch_start_time_point
	if backtest_start_timepoint.Before(fetch_start_timepoint) {
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
func get_price(db *sql.DB, assets []asset.Stock, asset_names []string, start_time time.Time, price_type int) error {
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

func fetch_price(db *sql.DB, asset_names []string, start_time time.Time, algo_mess []Algo_Message) error {
	for i, v := range asset_names {
		var err error
		algo_mess[i].Asset_Name = v
		algo_mess[i].P, err = asset.Read_Next_Data(db, v, start_time)
		if err != nil {
			return err
		}
	}
	return nil
}

// goroutine get each price algo need, and pass to algo goroutine via channel
func fetcher(db *sql.DB, fetcher_init_chan chan Fetcher_Init, messages chan []Algo_Message, ter_chan chan Algo_Terminal_And_Statistical, Err_Chan chan Err_Message) {
	//get init message from algo runner
	init_message := <-fetcher_init_chan
	//for send to algo runner
	algo_mess := make([]Algo_Message, len(init_message.Asset_Names))
	start_time := init_message.Start_TimePoint

	for {
		err := fetch_price(db, init_message.Asset_Names, start_time, algo_mess)
		if err != nil {
			//fetcher发生错误将信息回传给main goroutine,然后退出此goroutine
			Err_Chan <- Err_Message{Err: err, Gorotuine_Type: "fetcher"}
			return
		}
		start_time = start_time.Add(24 * time.Hour)
		//向algo runner发送message，通过Algo_Message channel，此时如果algo runner没有从channel中拿数据就会被阻塞
		messages <- algo_mess
		select {
		case <-ter_chan:
			return
			//因为我们不希望一直被阻塞在<-ter_chan这里,所以用time.After，在执行到ter_chan被阻塞300毫秒(0.3秒)后没有从ter_chan拿到ter_chan消息就break，fetcher取下一个数据
		case <-time.After(300 * time.Millisecond):
			break
		}
	}

}

func algo_runner(algo_init_chan chan AlgoRunner_Init, fetcher_init_chan chan Fetcher_Init, algo_mess_chan chan []Algo_Message, stat_and_ter_chan chan Algo_Terminal_And_Statistical, err_mess_chan chan Err_Message) error {
	//从channel AlgoRunner_Init中拿到algorunner的初始化数据
	algo_init_mess := <-algo_init_chan

	switch {
	case algo_init_mess.Algo == "Stat_Arb":
		params := algolib.Params{}
		params.IsBackTest = true
		params.S = algo_init_mess.Assets
		params.Backtest_Start_Time = algo_init_mess.Backtest_Start_TimePoint
		algolib.Call_Algo(params, algolib.Stat_Arb)

	case algo_init_mess.Algo == "Mean_Reversion":

	default:
		return errors.New("Algo Input Error: Algo Not Found")
	}
	return nil
}

func statistical_of_backtest() {

}

/*
这里讲述一个backtest的流程思路,首先有2个worker goroutine分别是
1: fetcher:不断地fetch数据给algo_runner
2: algo_runner:拿到fetcher的数据然后处理
首先是main传递初始化信息给algo_runner,使algo_runner进行初始化

流程见本目录的backtest_architecture_design.vsdx
*/
func Backtest_Main(db *sql.DB) error {
	//用于main goroutine传递初始化给algo runnter goroutine
	AlgoRunner_Init_Chan := make(chan AlgoRunner_Init)
	//用于algo runner发送初始化信息给Fetcher goroutine
	Fetcher_Init_Chan := make(chan Fetcher_Init)
	//用于Fetcher发送message给algo runner
	Algo_message_Chan := make(chan []Algo_Message)
	//用于algo runner goroutine发送terminal信号给Fetcher goroutine，和Statistical给main
	Algo_Ter_Stat_Chan := make(chan Algo_Terminal_And_Statistical)
	//用于fetcher和algo_runner回传错误信息给main goroutine
	Err_Chan := make(chan Err_Message)

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
	var fetch_start_timepoint_raw string
	var fetch_start_timepoint time.Time
	fetch_start_timepoint_raw = os.Getenv("FETCH_START_TIMEPOINT")
	err = get_fetch_start_timepoint(fetch_start_timepoint_raw, &fetch_start_timepoint)
	if err != nil {
		return err

	}

	//从.env中获取backtest_start_time_point
	var backtest_start_timepoint_raw string
	var backtest_start_timepoint time.Time
	backtest_start_timepoint_raw = os.Getenv("BACKTEST_START_TIMEPOINT")
	//无论是否从.env获取到了BACKTEST_START_TIMEPOINT，都要进入下面函数进行判断
	err = get_backtest_start_timepoint(backtest_start_timepoint_raw, &backtest_start_timepoint, fetch_start_timepoint)
	if err != nil {
		return err
	}

	//get price
	var assets []asset.Stock
	err = get_price(db, assets, asset_names, fetch_start_timepoint, price_type)
	if err != nil {
		return err
	}

	//运行fetcher
	go fetcher(db, Fetcher_Init_Chan, Algo_message_Chan, Algo_Ter_Stat_Chan, Err_Chan)

	//运行algo_runner
	go algo_runner(AlgoRunner_Init_Chan, Fetcher_Init_Chan, Algo_message_Chan, Algo_Ter_Stat_Chan, Err_Chan)

	AlgoRunner_Init_Chan <- AlgoRunner_Init{
		Algo:                     algo,
		Backtest_Start_TimePoint: backtest_start_timepoint,
		Assets:                   assets,
	}
	//statistical of backtest

	return nil

}
