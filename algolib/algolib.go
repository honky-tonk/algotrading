package algolib

import (
	"algotrading/asset"
	"algotrading/indicator"
)

// callback func of trading algo
type Trading_Algo func(params Params) (Statistc_Result, error)

type Params struct {
	IsBackTest bool
	Factors    indicator.Factors
	S          []asset.Stock
}

type Message struct {
	S asset.Stock
}

// 回测的结果统计信息
type Statistc_Result struct {
}

// use callback func of trading algo, param 1 is param of callback func
func Call_Algo(param Params, algo Trading_Algo) {
	//call trading algo
	algo(param)
}
