package asset

import (
	"algotrading/global"
	_ "algotrading/global"
	"algotrading/logger"
	"database/sql"
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

// for daily price
// ptype is price type, sname is stock name, we will fill price map
func get_daily_price(ptype string, sname string, start_timepoint time.Time, db *sql.DB) ([]Price, error) {
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
	if len(d.Time_Series) == 0 {
		return nil, errors.New("get price of asset from alphavantage error!")
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

	start_index := len(s) - int(time.Now().Sub(start_timepoint).Hours()/24)
	s = s[start_index:]
	for i, p := range s {
		if p.T == start_timepoint {
			start_index = i
		}
	}
	s = s[start_index:]

	//write to database
	err = Write_To_Database(db, sname, s)
	if err != nil {
		logger.Info.Println("write database error: " + err.Error())
		return nil, err
	}
	return s, nil

}

// get price pre weekly(friday night price)
func get_weekly_price(db *sql.DB, ptype string, sname string, start_timepoint time.Time) ([]Price, error) {
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
	if len(w.Time_Series) == 0 {
		return nil, errors.New("get price of asset from alphavantage error!")
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
	//write to database
	err = Write_To_Database(db, sname, s)
	if err != nil {
		logger.Info.Println("write database error: " + err.Error())
		return nil, err
	}
	return s, nil

}

// get the price of each month last trade day's
func get_monthly_price(db *sql.DB, ptype string, sname string, period int) ([]Price, error) {
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
	if len(m.Time_Series) == 0 {
		return nil, errors.New("get price of asset from alphavantage error!")
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
	//write to database
	err = Write_To_Database(db, sname, s)
	if err != nil {
		logger.Info.Println("write database error: " + err.Error())
		return nil, err
	}
	return s, nil
}

func need_update_data(p []Price) bool {
	now := time.Now()          //time of now
	db_newest := p[len(p)-1].T //time of newst data in database

	oneday, _ := time.ParseDuration("24h")
	//threeday, _ := time.ParseDuration("72h")

	fmt.Println("now is: ", now.String(), "db_newest is: ", db_newest.String())
	if now == db_newest || now.Sub(db_newest) <= oneday { //with one day error
		return false //no need update
	}

	return true
}

func (s *Stock) Check_Stock_Exist_From_Database(db *sql.DB) bool {
	var exist bool

	query := `SELECT EXISTS(SELECT * FROM sh_stock  WHERE stock_id = $1 AND time = $2);`
	rows, err := db.Query(query, s.Name, s.Start_TimePoint)
	if err != nil {
		logger.Error.Fatal("SQL can't exec:", err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&exist)
		if err != nil {
			logger.Error.Fatal("can't scan result: ", err)
		}
	}

	return exist

}

func Write_To_Database(db *sql.DB, sname string, s []Price) error {
	//if exist then not insert record
	tx, _ := db.Begin()
	query := `SELECT EXISTS(SELECT 1 FROM sh_stock WHERE stock_id = $1 AND  time = $2);`
	insert := `INSERT INTO sh_stock (stock_id, time, open, close, high, low, volume) VALUES ($1, $2, $3, $4, $5, $6, $7);`
	for _, i := range s {
		var exist bool
		row := tx.QueryRow(query, sname, i.T.Format("2006-01-02"))
		err := row.Scan(&exist)
		if err != nil {
			tx.Rollback()
			logger.Info.Println("can't scan result: ", err)
			return errors.New(err.Error())
		}
		if exist == false {
			_, err := tx.Exec(insert, sname, i.T.Format("2006-01-02"), i.SP.Open, i.SP.Close, i.SP.High, i.SP.Low, i.SP.Volume)
			if err != nil {
				tx.Rollback()
				logger.Info.Println("can't exec sql: ", err)
				return errors.New("can't exec sql: " + err.Error())
			}
		}
	}
	err := tx.Commit()
	if err != nil {
		tx.Rollback()
		logger.Info.Println("can't rollback  ", err)
		return errors.New("can't rollback  " + err.Error())
	}
	return nil
}

// 给定一个时间读下一条数据
func Read_Next_Data(db *sql.DB, sname string, start_timepoint time.Time) (Price, error) {
	//read database first
	query := `SELECT * FROM sh_stock WHERE stock_id = $1 AND time > $2 ORDER BY time LIMIT 1;`
	//Price for return
	var p Price
	//tmp str
	tmp_str := "tmp"
	//start tx
	tx, _ := db.Begin()
	//query
	row := db.QueryRow(query, sname, start_timepoint)
	//get query result
	err := row.Scan(&tmp_str, &p.T, &p.SP.Open, &p.SP.Close, &p.SP.High, &p.SP.Low, &p.SP.Volume)
	if err != nil {
		logger.Info.Println("sql result can't scanf : ", err.Error())
		tx.Rollback()
		return Price{}, err
	}
	//commit
	err = tx.Commit()
	if err != nil {
		logger.Info.Println("sql can't commit : ", err.Error())
		tx.Rollback()
		return Price{}, err
	}
	//read success from database(row exist in database)
	if tmp_str != "tmp" {
		return p, nil
	}

	//read from database failed(row not exist in database), then read from network

}

func Read_Stock_Data_From_Database(d *sql.DB, sname string, start_timepoint time.Time) ([]Price, error) {
	var p []Price //for return
	var tmp_price Price
	query := `SELECT * FROM sh_stock WHERE stock_id = $1 AND time > $2 ORDER BY time DESC;`
	tx, _ := d.Begin()
	//query
	rows, err := tx.Query(query, sname, start_timepoint)
	if err != nil {
		logger.Info.Println("can't exec sql: ", err.Error())
		tx.Rollback()
		return nil, err
	}
	//fill price for return
	for rows.Next() {
		var temp string
		err := rows.Scan(&temp, &tmp_price.T, &tmp_price.SP.Open, &tmp_price.SP.Close, &tmp_price.SP.High, &tmp_price.SP.Low, &tmp_price.SP.Volume)
		//fmt.Println("v is : ", v)
		p = append(p, tmp_price)
		if err != nil {
			logger.Info.Println("sql result can't scanf : ", err.Error())
			tx.Rollback()
			return nil, err
		}
	}
	//commit
	err = tx.Commit()
	if err != nil {
		logger.Info.Println("sql can't commit : ", err.Error())
		tx.Rollback()
		return nil, err
	}

	//fmt.Println(p)
	return p, nil
}

func (s *Stock) Get_Price(d *sql.DB) (err error) {
	//var time_type string
	//var price_from_api interface{}

	fmt.Println("Daily is: ", global.Daily)
	switch {
	case s.Type == global.Daily:
		//check if exist in database
		exist := s.Check_Stock_Exist_From_Database(d)
		if exist == true {
			//read data from local database
			s.Prices, err = Read_Stock_Data_From_Database(d, s.Name, s.Start_TimePoint)
			if err != nil {
				return err
			}
			//fmt.Println(s.Prices)
			break
		}
		s.Prices, err = get_daily_price("TIME_SERIES_DAILY", s.Name, s.Start_TimePoint, d)
		if err != nil {
			return err
		}
		//fmt.Println(s.Prices)
		break
	case s.Type == global.Weekly:
		//check if exist in database
		exist := s.Check_Stock_Exist_From_Database(d)
		if exist == true {
			//read data from local database
			s.Prices, err = Read_Stock_Data_From_Database(d, s.Name, s.Start_TimePoint)
			if err != nil {
				return err
			}
			break
		}
		//get from alphavantage
		s.Prices, err = get_weekly_price(d, "TIME_SERIES_WEEKLY", s.Name, s.Start_TimePoint)
		if err != nil {
			return err
		}
		break
	case s.Type == global.Monthly:
		//check if exist in database
		exist := s.Check_Stock_Exist_From_Database(d)
		if exist == true {
			//read data from local database
			s.Prices, err = Read_Stock_Data_From_Database(d, s.Name, s.Start_TimePoint)
			if err != nil {
				return err
			}
			break
		}
		//get from alphavantage
		s.Prices, err = get_monthly_price(d, "TIME_SERIES_MONTHLY", s.Name, s.Start_TimePoint)
		if err != nil {
			return err
		}
		break
	default:
		return errors.New("error stock time type")
	}

	return nil
}
