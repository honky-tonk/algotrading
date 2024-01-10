package indicator

import (
	"algotrading/asset"
	"errors"
)

func (m MACD_Indicator) Set_Period(period int) { /*macd no need period*/ }

/*
MACD = 12_period EMA - 26_period EMA
*/
func (m *MACD_Indicator) Calculate_Indicator(s *asset.Stock) error {
	m.Smoothing_EMA = 2

	ema_12_period_indic := EMA_Indicator{}
	ema_26_period_indic := EMA_Indicator{}
	signal_indic := EMA_Indicator{}

	if m.Smoothing_EMA == 0 {
		return errors.New("Please fill smoothing member of EMA_Indicator struct obj")
	}

	if len(s.Prices) < 26 {
		return errors.New("Not Have enough sample for indicator")
	}

	//init 12 period of ema indicator
	//ema_12_period_indic.Asset_Type = s.Type
	ema_12_period_indic.Smoothing = m.Smoothing_EMA
	ema_12_period_indic.Period = 12

	//init 26 period of ema indicator
	//ema_26_period_indic.Asset_Type = s.Type
	ema_26_period_indic.Smoothing = m.Smoothing_EMA
	ema_26_period_indic.Period = 26

	//init signal indicator
	//signal_indic.Asset_Type = s.Type
	signal_indic.Smoothing = m.Smoothing_EMA
	signal_indic.Period = 9

	var err error
	//get indicitor of 12_period of ema
	err = ema_12_period_indic.Calculate_Indicator(s)
	if err != nil {
		return err
	}
	//get indicator of 26_period of ema
	err = ema_26_period_indic.Calculate_Indicator(s)
	if err != nil {
		return err
	}

	//get macd indicate
	offset := ema_26_period_indic.Period - ema_12_period_indic.Period
	for i := 0; i < len(ema_26_period_indic.Indicator_Value); i++ {
		tmp_p := asset.Indicator_Value{}
		if ema_26_period_indic.Indicator_Value[i].T == ema_12_period_indic.Indicator_Value[i+offset].T {
			tmp_p.P = ema_12_period_indic.Indicator_Value[i+offset].P - ema_26_period_indic.Indicator_Value[i].P
			tmp_p.T = ema_26_period_indic.Indicator_Value[i].T
			m.MACD_Indicator_Value = append(m.MACD_Indicator_Value, tmp_p)
		}
	}

	//get 9_period ema indicator for signal
	m.Signal_Indicator, err = signal_indic.Calculate_Indicator_macd_sig(m.MACD_Indicator_Value)
	if err != nil {
		return err
	}

	return nil
}
