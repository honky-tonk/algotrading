package backtest

import (
	"time"

	"algotrading/asset"
	"algotrading/indicator"
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
	Asset_Type      int
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
	Err        Err_Message
}

type Err_Message struct {
	Gorotuine_Type string
	Err            error
}

// algo_runner运行完后的statistical信息
type Statistical struct {
	Max_DrawDown float64
	Sharp_Ratio  float64
	Start_Price  float64
	Return       float64
}

// callback func of trading algo
type Trading_Algo func(params Params)

type Params struct {
	//for algo
	IsBackTest          bool
	Factors             indicator.Factors
	S                   []asset.Stock
	Backtest_Start_Time time.Time

	//use for backtest platform
	Algo_Init_Chan    chan AlgoRunner_Init
	Fetcher_Init_Chan chan Fetcher_Init
	Algo_Mess_Chan    chan []Algo_Message
	Ter_Stat_Chan     chan Algo_Terminal_And_Statistical
	//Err_Mess_Chan     chan Err_Message
}

type Message struct {
	S asset.Stock
}

// //Stat Arb
// //每一次进行套利计算后得到的信息
// type Stat_Arb_Info struct{
// 	float64

// }

// use callback func of trading algo, param 1 is param of callback func
func Call_Algo(param Params, algo Trading_Algo) {
	//call trading algo
	algo(param)
}
