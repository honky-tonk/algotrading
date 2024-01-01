package algolib

import (
	"algotrading/asset"
	"errors"
	"time"
)

// cause assets time of price not match so correct, delete not match price of data
func Correct_Price(s []asset.Stocks) ([]asset.Stocks, error) {
	if len(s) != 2 {
		return nil, errors.New("Correct Price error!: Please input two slice of Stocks")
	}
	//for result
	result := make([]asset.Stocks, 2)

	//create map
	asset1 := make(map[time.Time]asset.Stock_Price)
	//store asset1 to map
	for _, v := range s[0].Prices {
		asset1[v.T] = v.SP
	}
	//compare with asset2
	for _, v := range s[1].Prices {

	}

	return result, nil
}

func Conver_Stocks_To_Float64Slices(s asset.Stocks) (string, []float64) {
	result := make([]float64, 0)
	for _, v := range s.Prices {
		result = append(result, v.SP.Close)
	}
	return s.Name, result
}

// find 2个asset是否是correlation的，只要correlation在0.5到0.9之间都称为correlation
// 输入多个assets,返回最为correlation的2个
func Find_Correlation(assets []asset.Stocks) ([]asset.Stocks, error) {
	var max_correlated float64
	//only compare two assets
	if len(assets) >= 2 {
		return nil, errors.New("Please Input More than two Assets")
	}
	//for return
	correlated_assets := make([]asset.Stocks, 2)
	//find max_correlated
	//TODO
	return correlated_assets, nil

}

// Stat_Arb执行前2个asset已经是correlation
func Stat_Arb(param Params) (Statistc_Result, error) {
	//result for return
	result := Statistc_Result{}
	correlated_assets, err := Find_Correlation(param.S)
	if err != nil {
		return result, err
	}
	//TODO
	return result, nil
}
