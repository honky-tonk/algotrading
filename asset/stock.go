package asset

import (
	"algotrading/global"
	_ "algotrading/global"
	"algotrading/logger"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	//"slices" slices.Reverse support since go 1.21
	"sort"
	"strconv"
	"time"
)

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

//for program
type Stock_Price struct {
	Open   float64 `json:"open"`
	High   float64 `json:"high"`
	Low    float64 `json:"low"`
	Close  float64 `json:"close"`
	Volume int64   `json:"volume"`
}

type Price struct {
	T  time.Time   `json:"time"`
	SP Stock_Price `json:"stock_price"`
}

type Stocks struct {
	Prices []Price
	Type   int    `json:"data_type"`
	Name   string `json:"stock_name"`
	Period int    `json:"period"`
}

func get_price_from_api(ptype string, assert_name string) (*http.Response, error) {
	url := fmt.Sprintf("%sfunction=%s&outputsize=full&symbol=%s&apikey=%s", global.Stock_Api, ptype, assert_name, global.Api_Key)
	logger.Info.Println("get url: ", url)
	return http.Get(url)
}

func reverse_slice(s []Price) []Price {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	return s
}

//for daily price
// ptype is price type, sname is stock name, we will fill price map
func get_daily_price(ptype string, sname string, period int) ([]Price, error) {
	d := Daily_Stock{}
	resp, err := get_price_from_api(ptype, sname)
	if err != nil {
		return nil, err
	}
	//read from respond body
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// unmarshal(take a serialized object to internal data structure) full Daily_Stock struct
	err = json.Unmarshal(b, &d)
	if err != nil {
		return nil, err
	}

	//convert map[string]Api_string which  get from internet to slice of Price
	s := make([]Price, len(d.Time_Series))
	i := 0
	for k, v := range d.Time_Series {
		time, err := time.Parse("2006-01-02", k)
		if err != nil {
			return nil, err
		}
		tmp_price := Stock_Price{}
		tmp_price.Close, err = strconv.ParseFloat(v.Close, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Open, err = strconv.ParseFloat(v.Open, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.High, err = strconv.ParseFloat(v.High, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Low, err = strconv.ParseFloat(v.Low, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Volume, err = strconv.ParseInt(v.Volume, 10, 64)
		if err != nil {
			return nil, err
		}

		s[i] = Price{T: time, SP: tmp_price}
		i++
	}
	//sort
	sort.Slice(s, func(i, j int) bool {
		return s[j].T.Before(s[i].T)
	})
	s = s[:period]
	s = reverse_slice(s)
	return s, nil

}

//get price pre weekly(friday night price)
func get_weekly_price(ptype string, sname string, period int) ([]Price, error) {
	w := Weekly_Stock{}
	resp, err := get_price_from_api(ptype, sname)
	if err != nil {
		return nil, err
	}
	//read from respond body
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	// unmarshal(take a serialized object to internal data structure) full Weekly_Stock struct
	err = json.Unmarshal(b, &w)
	if err != nil {
		return nil, err
	}

	//convert map[string]Api_string which  get from internet to slice of Price
	s := make([]Price, len(w.Time_Series))
	i := 0
	for k, v := range w.Time_Series {
		time, err := time.Parse("2006-01-02", k)
		if err != nil {
			return nil, err
		}
		tmp_price := Stock_Price{}
		tmp_price.Close, err = strconv.ParseFloat(v.Close, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Open, err = strconv.ParseFloat(v.Open, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.High, err = strconv.ParseFloat(v.High, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Low, err = strconv.ParseFloat(v.Low, 64)
		if err != nil {
			return nil, err
		}

		tmp_price.Volume, err = strconv.ParseInt(v.Volume, 10, 64)
		if err != nil {
			return nil, err
		}

		s[i] = Price{T: time, SP: tmp_price}
		i++
	}
	//sort
	sort.Slice(s, func(i, j int) bool {
		return s[j].T.Before(s[i].T)
	})
	s = s[:period]
	s = reverse_slice(s)
	return s, nil

}

//get the price of each month last trade day's
func get_monthly_price(ptype string, sname string, period int) ([]Price, error) {
	m := Monthly_Stock{}
	resp, err := get_price_from_api(ptype, sname)
	if err != nil {
		return nil, err
	}
	//read from respond body
	b, err := io.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	//fmt.Println(b)
	// unmarshal(take a serialized object to internal data structure) full Monthly_Stock struct
	err = json.Unmarshal(b, &m)
	//fmt.Println(d.Meta_Datas)
	if err != nil {
		//fmt.Println(err.Error())
		return nil, err
	}

	//convert map[string]Api_string which  get from internet to slice of Price
	s := make([]Price, len(m.Time_Series))
	i := 0
	for k, v := range m.Time_Series {
		//fmt.Println("key: ", k, "value: ", v)
		time, err := time.Parse("2006-01-02", k)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}
		tmp_price := Stock_Price{}
		tmp_price.Close, err = strconv.ParseFloat(v.Close, 64)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}

		tmp_price.Open, err = strconv.ParseFloat(v.Open, 64)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}

		tmp_price.High, err = strconv.ParseFloat(v.High, 64)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}

		tmp_price.Low, err = strconv.ParseFloat(v.Low, 64)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}

		tmp_price.Volume, err = strconv.ParseInt(v.Volume, 10, 64)
		if err != nil {
			//fmt.Println(err.Error())
			return nil, err
		}

		//s[i] = price_entries{time: tmp_price}
		s[i] = Price{T: time, SP: tmp_price}
		i++
	}
	//sort
	sort.Slice(s, func(i, j int) bool {
		return s[j].T.Before(s[i].T)
	})
	s = s[:period]
	s = reverse_slice(s)
	return s, nil
}

func (s *Stocks) Get_Price() (err error) {
	//var time_type string
	//var price_from_api interface{}

	fmt.Println("Daily is: ", Daily)
	switch {
	case s.Type == Daily:
		s.Prices, err = get_daily_price("TIME_SERIES_DAILY", s.Name, s.Period)
		if err != nil {
			return err
		}
		break
	case s.Type == Weekly:
		s.Prices, err = get_weekly_price("TIME_SERIES_WEEKLY", s.Name, s.Period)
		if err != nil {
			return err
		}
		break
	case s.Type == Monthly:
		s.Prices, err = get_monthly_price("TIME_SERIES_MONTHLY", s.Name, s.Period)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("error stock time type")
	}

	return nil
}
