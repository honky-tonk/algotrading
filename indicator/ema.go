package indicator

import (
	"algotrading/asset"
	"errors"
	"fmt"
	//"fmt"
)

/*
we use SMA for init ema init indicator
*/

func (ema *EMA_Indicator) Set_Period(period int) { ema.Period = period }

func (ema *EMA_Indicator) Calculate_Indicator(s *asset.Stock) error {
	ema.Smoothing = 2
	if ema.Smoothing == 0 {
		return errors.New("Please fill smoothing member of EMA_Indicator struct obj")
	}

	if ema.Period == 0 {
		return errors.New("Please set Period member of EMA_Indicator struct obj")
	}

	if len(s.Prices) < ema.Period {
		return errors.New("Not Have enough sample for indicator")
	}

	var Multipler float64
	Multipler = float64(ema.Smoothing) / float64((1 + ema.Period))

	//for init ema indicator
	prices := make([]float64, 0)
	for i := 0; i < ema.Period; i++ {
		prices = append(prices, s.Prices[i].SP.Close)
	}

	init_ema_price := sum_of_slice(prices) / float64(ema.Period)
	//fmt.Printf("get init_ema_price %f, multipler is %f\n", init_ema_price, Multipler)
	ema.Indicator_Value = append(ema.Indicator_Value, asset.Indicator_Value{P: init_ema_price, T: s.Prices[ema.Period-1].T})

	/*EMA_today = Value_today.(Smoothing/(1+Period)) + EMA_yesterday.(1 - (Smoothing/(1+Period)))*/
	for i := ema.Period; i < len(s.Prices); i++ {
		yeasterday_ema := ema.Indicator_Value[(len(ema.Indicator_Value) - 1)].P
		today_price := s.Prices[i].SP.Close
		ema.Indicator_Value = append(ema.Indicator_Value, asset.Indicator_Value{P: today_price*Multipler + (yeasterday_ema * (1 - Multipler)), T: s.Prices[i].T})
		//fmt.Println("yesterday is ", yeasterday_ema, "today's price is ", today_price, " ", ema.Indicator_Value[len(ema.Indicator_Value)-1], " ", s.Prices[i].T)
	}
	fmt.Println("------------------debug--------------------\n", ema.Indicator_Value, "----------------debug---------------\n")
	return nil
}

func (sig EMA_Indicator) Calculate_Indicator_macd_sig(p []asset.Indicator_Value) ([]asset.Indicator_Value, error) {
	if sig.Smoothing == 0 {
		return nil, errors.New("Please fill smoothing member of EMA_Indicator struct obj")
	}

	if sig.Period == 0 {
		return nil, errors.New("Please set Period member of EMA_Indicator struct obj")
	}

	if len(p) < sig.Period {
		return nil, errors.New("Not Have enough sample for indicator")
	}

	var Multipler float64
	Multipler = float64(sig.Smoothing) / float64((1 + sig.Period))

	//for init ema indicator
	prices := make([]float64, 0)
	for i := 0; i < sig.Period; i++ {
		prices = append(prices, p[i].P)
	}

	init_ema_price := sum_of_slice(prices) / float64(sig.Period)
	//fmt.Printf("get init_ema_price %f, multipler is %f\n", init_ema_price, Multipler)
	sig.Indicator_Value = append(sig.Indicator_Value, asset.Indicator_Value{P: init_ema_price, T: p[sig.Period-1].T})

	/*EMA_today = Value_today.(Smoothing/(1+Period)) + EMA_yesterday.(1 - (Smoothing/(1+Period)))*/
	for i := sig.Period; i < len(p); i++ {
		yeasterday_ema := sig.Indicator_Value[(len(sig.Indicator_Value) - 1)].P
		today_price := p[i].P
		sig.Indicator_Value = append(sig.Indicator_Value, asset.Indicator_Value{P: today_price*Multipler + (yeasterday_ema * (1 - Multipler)), T: p[i].T})
		//fmt.Println("yesterday is ", yeasterday_ema, "today's price is ", today_price, " ", ema.Indicator_Value[len(ema.Indicator_Value)-1], " ", s.Prices[i].T)
	}

	return sig.Indicator_Value, nil
}
