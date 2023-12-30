package algolib

import (
	"algotrading/asset"
	"errors"
)

// find 2个asset是否是correlation的，只要correlation在0.5到0.9之间都称为correlation
// 输入多个assets,返回最为correlation的2个
func Find_Correlation(assets []asset.Stocks) ([]asset.Stocks, error) {
	//only compare two assets
	if len(assets) >= 2 {
		return nil, errors.New("Please Input More than two Assets")
	}
	//for return
	correlated_assets := make([]asset.Stocks, 2)
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
