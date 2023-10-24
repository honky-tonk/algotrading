package asset

import (
	"algotrading/global"
	_ "algotrading/global"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
)

// for api
type Daily_Stock struct {
	Meta_Data   string               `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Time Series (Daily)"`
}

type Weekly_Stock struct {
	Meta_Data   string               `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Weekly Time Series"`
}

type Monthly_Stock struct {
	Meta_Data   string               `json:"Meta Data"`
	Time_Series map[string]Api_Price `json:"Monthly Time Series"`
}

type Api_Price struct {
	Open   string `json:"1. open"`
	Close  string `json:"2. high"`
	High   string `json:"3. low"`
	Low    string `json:"4. close"`
	Volume string `json:"5. volume"`
}

//for program
type Stock_Price struct {
	Open   float64
	Close  float64
	High   float64
	Low    float64
	Volume int
}

type Stocks struct {
	Price map[string]Stock_Price
	Type  int
	Name  string
}

func (s *Stocks) Get_Daily_Price(period int) (err error) {
	var time_type string
	var price_from_api interface{}

	switch {
	case s.Type == Daily:
		time_type = "TIME_SERIES_DAILY"
		price_from_api = new(Daily_Stock)
		break
	case s.Type == Weekly:
		time_type = "TIME_SERIES_WEEKLY"
		price_from_api = new(Weekly_Stock)
		break
	case s.Type == Monthly:
		time_type = "TIME_SERIES_MONTHLY"
		price_from_api = new(Monthly_Stock)
		break
	default:
		return errors.New("error stock time type")
	}

	url := fmt.Sprintf("%sfunction=%s&outputsize=full&symbol=%s&apikey=%s", global.Stock_Api, time_type, s.Name, global.Api_Key)
	fmt.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	json.Unmarshal(body, price_from_api)
	//fmt.Println(price_from_api)

	slice := reflect.ValueOf(price_from_api)

	for i := slice.Len(); i > period; i-- {
		fmt.Println(slice.Index(i - 1))
	}

	return nil
}
