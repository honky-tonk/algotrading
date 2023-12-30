package algolib

import (
	"algotrading/asset"
	"algotrading/indicator"
)

// callback func of trading algo
type Trading_Algo func(params Params)

type Params struct {
	Factors indicator.Factors
	S       []asset.Stocks
}

// use callback func of trading algo, param 1 is param of callback func
func Call_Algo(param Params, algo Trading_Algo) {
	//call trading algo
	algo(param)
}
