package backtest

import (
	"algotrading/asset"
	"errors"
	"time"

	"gonum.org/v1/gonum/stat"
)

// is slice contains t
func is_cotains(s []asset.Price, t time.Time) bool {
	for _, v := range s {
		if v.T.Equal(t) {
			return true
		}
	}
	return false
}

// cause assets time of price not match so correct, delete not match price of data
func Correct_Price(s1 asset.Stock, s2 asset.Stock) (asset.Stock, asset.Stock) {
	//for result
	result := make([]asset.Stock, 2)
	result[0] = s1
	result[1] = s2

	//delete asset.Stock1 element which exist asset.Stock1, but not in asset.Stock2
	for i := 0; i < len(result[0].Prices); i++ {
		if !is_cotains(result[1].Prices, result[0].Prices[i].T) {
			result[0].Prices = append(result[0].Prices[:i], result[0].Prices[i+1:]...)
		}
	}

	//delete asset.Stock2 element which exist asset.Stock2, but not in asset.Stock1
	for i := 0; i < len(result[1].Prices); i++ {
		if !is_cotains(result[0].Prices, result[1].Prices[i].T) {
			result[1].Prices = append(result[1].Prices[:i], result[1].Prices[i+1:]...)
		}
	}

	return result[0], result[1]
}

func Conver_Stocks_To_Float64Slices(s asset.Stock) []float64 {
	result := make([]float64, 0)
	for _, v := range s.Prices {
		result = append(result, v.SP.Close)
	}
	return result
}

// find 2个asset是否是correlation的，只要correlation在0.5到0.9之间都称为correlation
// 输入多个assets,返回最为correlation的2个
func Find_Max_Correlation(assets []asset.Stock) (asset.Stock, asset.Stock, float64, error) {
	//for return
	var max_correlated float64
	var max_correlated_index1 int
	var max_correlated_index2 int

	//only compare two assets
	if len(assets) >= 2 {
		return asset.Stock{}, asset.Stock{}, 0.0, errors.New("Please Input More than two Assets")
	}
	max_correlated = 0.0
	max_correlated_index1 = -1
	max_correlated_index2 = -1

	for i := 0; i < len(assets); i++ {
		for j := i + 1; j < len(assets); i++ {
			s1, s2 := Correct_Price(assets[i], assets[j])
			prices1 := Conver_Stocks_To_Float64Slices(s1)
			prices2 := Conver_Stocks_To_Float64Slices(s2)
			corr := stat.Correlation(prices1, prices2, nil)
			if corr > max_correlated && (corr >= 0.5 && corr <= 0.9) {
				max_correlated = corr
				max_correlated_index1 = i
				max_correlated_index2 = j
			}
		}
	}

	//not found corr
	if max_correlated == 0.0 {
		return asset.Stock{}, asset.Stock{}, 0.0, nil
	}

	return assets[max_correlated_index1], assets[max_correlated_index2], max_correlated, nil

}

// Stat_Arb执行前2个asset已经是correlation
func Stat_Arb(param Params) (Statistc_Result, error) {
	//result for return
	result := Statistc_Result{}

	//find max correalation
	asset1, asset2, corr, err := Find_Max_Correlation(param.S)
	if err != nil {
		return result, err
	}
	if corr == 0.0 {
		//return not found correlation assets
	}

	//TODO
	return result, nil
}
