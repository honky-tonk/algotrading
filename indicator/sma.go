package indicator

/*for Simple Moving Average(SMA) indicator*/
import (
	"algotrading/asset"
	"errors"
	"fmt"
)

type SMA_Indicator struct {
	//Asset_Type      int
	Period          int                     `json:"indic_period"`
	Indicator_Value []asset.Indicator_Value `json:"indic_values"`
}

/*
a array of price of end of day is: 1(end of day1), 2(end of day2), 3(end of day3), 4(end of day4), 5(end of day5), 6, 7, 8, 9(end of day9)
sma period is 5 day
sma indicator is: nil(end of day1), nil(end of day2), nil(end of day3), nil(end of day4), (1+2+3+4+5)/5=3(end of day5), 2+3+4+5+6/5=2.8(end of day6)......(5+6+7+8+9)/5=7(end of day9)
*/

func (sma *SMA_Indicator) Set_Period(period int) {
	fmt.Println("period of sma is ", period)
	sma.Period = period
}

func (sma *SMA_Indicator) Calculate_Indicator(s *asset.Stocks) error {
	if len(s.Prices) <= sma.Period {
		return errors.New("Not Have enough sample for indicator")
	}

	values := make([]asset.Indicator_Value, 0)

	var tmp_prices []float64
	//inital first sma.period
	fmt.Println("sma period is ", sma.Period)
	for i := 0; i < sma.Period-1; i++ {
		tmp_close_price := s.Prices[i].SP.Close
		tmp_prices = append(tmp_prices, tmp_close_price)
	}
	//full rest of all sma indicator
	for i := sma.Period - 1; i < len(s.Prices); i++ {
		tmp_close_price := s.Prices[i].SP.Close
		tmp_prices = append(tmp_prices, tmp_close_price)
		values = append(values, asset.Indicator_Value{P: sum_of_slice(tmp_prices) / float64(sma.Period), T: s.Prices[i].T})
		tmp_prices = tmp_prices[1:]
	}

	sma.Indicator_Value = values
	return nil
}

func (sma SMA_Indicator) Calculate_Indicator_For_kdj(s []asset.Indicator_Value) ([]asset.Indicator_Value, error) {
	if len(s) <= sma.Period {
		return nil, errors.New("Not Have enough sample for indicator")
	}

	values := make([]asset.Indicator_Value, 0)

	var tmp_prices []float64
	//inital first sma.period
	for i := 0; i < sma.Period-1; i++ {
		tmp_close_price := s[i].P
		tmp_prices = append(tmp_prices, tmp_close_price)
	}
	//full rest of all sma indicator
	for i := sma.Period - 1; i < len(s); i++ {
		tmp_close_price := s[i].P
		tmp_prices = append(tmp_prices, tmp_close_price)
		values = append(values, asset.Indicator_Value{P: sum_of_slice(tmp_prices) / float64(sma.Period), T: s[i].T})
		tmp_prices = tmp_prices[1:]
	}

	return values, nil
}
