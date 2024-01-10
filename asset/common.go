package asset

import (
	_ "algotrading/global"
	"time"
)

// for backtest message
type Backtest_Message struct {
	Asset_Name string
	Price      Price
}

// for backtest message send one time
type Backtest_Messages struct {
	Message  Backtest_Message
	Mess_Num int
}

// for api
type Daily_Stock struct {
	Meta_Datas  Meta_Data            `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Time Series (Daily)"`
}

type Weekly_Stock struct {
	Meta_Datas  Meta_Data            `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Weekly Time Series"`
}

type Monthly_Stock struct {
	Meta_Datas  Meta_Data            `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Monthly Time Series"`
}

type Meta_Data struct {
	Info           string `json:"1. Information"`
	Symbol         string `json:"2. Symbol"`
	Last_Refreshed string `json:"3. Last Refreshed"`
	Output_Size    string `json:"4. Output Size"`
	Time_Zone      string `json:"5. Time Zone"`
}

type Api_Price struct {
	Open   string `json:"1. open"`
	High   string `json:"2. high"`
	Low    string `json:"3. low"`
	Close  string `json:"4. close"`
	Volume string `json:"5. volume"`
}

type Stock_Price struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

// for indicator
type Indicator_Value struct {
	T time.Time `json:"time"`
	P float64   `json:"value"`
}

type Price struct {
	T  time.Time   `json:"time"`
	SP Stock_Price `json:"stock_price"`
}

// for program
type Stock struct {
	Prices           []Price   `json:"prices"`
	Type             int       `json:"data_type"`
	Name             string    `json:"stock_name"`
	Start_TimePoint  time.Time `json:"start_timepoint"`
	Indicator_Type   int       `json:"indic"`
	Indicator_Period int       `json:"indic_period"`
}

type Asset interface {
	Get_Price(int) error
}
