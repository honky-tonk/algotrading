package indicator

import "algotrading/asset"

type Indicator interface {
	Calculate_Indicator(asset.Stocks) error
}

func sum_of_slice(s []float64) float64 {
	var r float64
	for _, v := range s {
		r += v
	}
	return r
}
