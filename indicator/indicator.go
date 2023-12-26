package indicator

import "algotrading/asset"

type Indicator interface {
	Calculate_Indicator(*asset.Stocks) error
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
