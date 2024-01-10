package indicator

import (
	"algotrading/asset"
	//"fmt"
)

/*for exponential moving average(EMA) indicator*/

type EMA_Indicator struct {
	//Asset_Type      int							`json:"type"`
	Period          int                     `json:"indic_period"`
	Indicator_Value []asset.Indicator_Value `json:"indic_values"`
	Smoothing       int                     `json:"smoothing"`
}

/*for KDJ indicator*/

type KDJ_Indicator struct {
	Kvalue []asset.Indicator_Value `json:"kvalues"`
	Dvalue []asset.Indicator_Value `json:"dvalues"`
	Jvalue []asset.Indicator_Value `json:"jvalues"`
	Period int                     `json:"indic_period"`
	//Type   int 						`json:"type"`
}

/*for moving average convergence/divergence indicator*/

type MACD_Indicator struct {
	//Asset_Type int
	//Period int # the formula of MACD is 12_Period EMA - 26_period EMA
	Signal_Indicator     []asset.Indicator_Value `json:"signal_values"` //9-period of ema(ema data from MACD)
	MACD_Indicator_Value []asset.Indicator_Value `json:"macd_values"`
	Smoothing_EMA        int                     `json:"smoothing_of_ema"`
	Period               int                     `json:"indic_period"`
}

/*for Simple Moving Average(SMA) indicator*/

type SMA_Indicator struct {
	//Asset_Type      int
	Period          int                     `json:"indic_period"`
	Indicator_Value []asset.Indicator_Value `json:"indic_values"`
}

type Indicator interface {
	Calculate_Indicator(*asset.Stock) error
	Set_Period(int)
}

// for backend return to front-end, if indicator_type != 4
type Stock_and_Indicator struct {
	Stock  []asset.Price `json:"prices"`
	Indic  Indicator     `json:"indic"`
	Period int           `json:"period"`
}

func sum_of_slice(s []float64) float64 {
	var r float64
	for _, v := range s {
		r += v
	}
	return r
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
