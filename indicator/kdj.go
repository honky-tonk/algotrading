package indicator

import (
	"algotrading/asset"
	"errors"
)

/*for KDJ indicator*/

type KDJ_Indicator struct {
	Kvalue []asset.Indicator_Value
	Dvalue []asset.Indicator_Value
	Jvalue []asset.Indicator_Value
	Type   int //asset type
}

// find max stock price give period and return INDIX
func find_max(p []asset.Price) int {
	max_value_idx := 0
	for i, _ := range p {
		if p[max_value_idx].SP.Close < p[i].SP.Close {
			max_value_idx = i
		}
	}

	return max_value_idx
}

// find min stock price give period and return INDIX
func find_min(p []asset.Price) int {
	min_value_idx := 0
	for i, _ := range p {
		if p[min_value_idx].SP.Close > p[i].SP.Close {
			min_value_idx = i
		}
	}

	return min_value_idx
}

// k value is (today close price - lowest price last n day) /(highest close price - lowest price last n day) where n is 14(of course you can edit this value)
func calculate_kvalue(s asset.Stocks) []asset.Indicator_Value {
	k := make([]asset.Indicator_Value, 0)
	var tmp_k_value asset.Indicator_Value
	for i, j := 0, 13; j < len(s.Prices); i, j = i+1, j+1 {

		max_value_idx := find_max(s.Prices[i:j+1]) + i //i is offset
		min_value_idx := find_min(s.Prices[i:j+1]) + i //i is offset

		tmp_k_value.P = (s.Prices[j].SP.Close - s.Prices[min_value_idx].SP.Close) / (s.Prices[max_value_idx].SP.Close - s.Prices[min_value_idx].SP.Close) * 100
		tmp_k_value.T = s.Prices[j].T

		k = append(k, tmp_k_value)
	}
	return k
}

// d value is SMA(3,kvalue)
func calculate_dvalue(t int, k []asset.Indicator_Value) ([]asset.Indicator_Value, int) {
	sma := SMA_Indicator{}
	sma.Asset_Type = t
	sma.Period = 3
	sma.Indicator_Value, _ = sma.Calculate_Indicator_For_kdj(k)

	return sma.Indicator_Value, sma.Period
}

// j value = 3 * %D - 2 * %k 3 and 2 is arbitrary scale
func calculate_jvalue(k []asset.Indicator_Value, d []asset.Indicator_Value, dvalue_period int) ([]asset.Indicator_Value, error) {
	j := make([]asset.Indicator_Value, 0)
	for i, _ := range d {
		if d[i].T != k[i+dvalue_period-1].T {
			return nil, errors.New("compare error, time not equal!!!")
		}

		tmp_value := asset.Indicator_Value{}
		tmp_value.P = float64(3)*d[i].P - float64(2)*k[i+dvalue_period-1].P
		tmp_value.T = d[i].T
		j = append(j, tmp_value)
	}
	return j, nil
}

func (kdj *KDJ_Indicator) Calculate_Indicator(s asset.Stocks) error {
	var kvalue []asset.Indicator_Value
	var dvalue []asset.Indicator_Value
	var jvalue []asset.Indicator_Value

	//k value period is 14 day, d value period of sma is 3
	if s.Period <= 17 {
		return errors.New("Period is not fit this indicator")
	}

	kvalue = calculate_kvalue(s)
	dvalue, period := calculate_dvalue(kdj.Type, kvalue)
	jvalue, err := calculate_jvalue(kvalue, dvalue, period)
	if err != nil {
		return err
	}

	kdj.Kvalue = kvalue
	kdj.Dvalue = dvalue
	kdj.Jvalue = jvalue
	return nil
}
