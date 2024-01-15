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

func get_spread(s1 asset.Stock, s2 asset.Stock) []float64 {
	var spread []float64
	//已经完成格式化直接对比
	for i, _ := range s1.Prices {
		spread = append(spread, s1.Prices[i].SP.Close-s2.Prices[i].SP.Close)
	}
	return spread
}

// Stat_Arb执行前2个asset已经是correlation
func Stat_Arb(param Params) {
	//terminal and stat message
	ter_stat := Algo_Terminal_And_Statistical{}
	//for fetcher init
	fetcher_init_mess := Fetcher_Init{}
	//get algo runner init message from main
	algo_init_mess := <-param.Algo_Init_Chan
	//get asset type
	var asset_type int
	//如果发现asset为空，说明从main拿到的algo_runner 初始化数据有问题直接退出程序
	if len(algo_init_mess.Assets) == 0 {
		ter_stat.Err.Gorotuine_Type = "fetcher"
		ter_stat.Err.Err = errors.New("Assets is null, error with init message")
		param.Ter_Stat_Chan <- ter_stat
		return
	}
	asset_type = algo_init_mess.Assets[0].Type
	//asset names
	var asset_names []string
	//init fetcher init message
	//fetcher_init_mess.Asset_Names = asset_names
	fetcher_init_mess.Asset_Type = asset_type
	fetcher_init_mess.Start_TimePoint = algo_init_mess.Backtest_Start_TimePoint
	//find max correalation
	asset1, asset2, corr, err := Find_Max_Correlation(param.S)
	if err != nil {
		ter_stat.Err.Gorotuine_Type = "fetcher"
		ter_stat.Err.Err = err
		param.Ter_Stat_Chan <- ter_stat
		return
	}
	if corr == 0.0 {
		ter_stat.Err.Gorotuine_Type = "fetcher"
		ter_stat.Err.Err = errors.New("Not found correlation assets")
		param.Ter_Stat_Chan <- ter_stat
		return
	}
	asset_names = append(asset_names, asset1.Name)
	asset_names = append(asset_names, asset2.Name)
	fetcher_init_mess.Asset_Names = asset_names
	//send to channel for fetcher init
	param.Fetcher_Init_Chan <- fetcher_init_mess

	spreads := get_spread(asset1, asset2)
	//回测
	for {
		select {
		//get from fetcher
		case mess := <-param.Algo_Mess_Chan:
			price1 := mess[0].P.SP.Close
			price2 := mess[1].P.SP.Close
			spread := price1 - price2
			spreads = append(spreads, spread)
			spread_sample_mean, spread_sample_stddev := stat.MeanStdDev(spreads, nil)
			spread_sample_zscore := (spread - spread_sample_mean) / spread_sample_stddev
			//TODO
		case <-time.After(time.Microsecond * 300):
			break
		//get err message from main or fetcher
		case <-param.Ter_Stat_Chan:
			return
		}
	}

	//param.Ter_Stat_Chan <- ter_stat
	//return

}
