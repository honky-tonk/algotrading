package handler

import (
	"algotrading/asset"
	"algotrading/db"
	"algotrading/indicator"
	"algotrading/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Asset_Query(c *gin.Context) {
	//open db
	d := db.Db_main()
	defer d.Close()

	s := asset.Stocks{}
	err := c.ShouldBindJSON(&s)
	if err != nil {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "please input correct format",
		})
		logger.Info.Println(err)
		return
	}

	fmt.Println("s.period is: ", s.Period)
	if s.Period == 0 || s.Name == "" || s.Type == 0 {
		c.JSON(400, gin.H{
			"code":    400,
			"message": "please full Period and name and type element",
		})
		logger.Info.Println(err)
		return
	}

	err = s.Get_Price(d)
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Server Internal Error",
		})
		logger.Info.Fatal(err)
		return
	}

	//for test indicator

	/*
		fmt.Println("--------------------price is -----------------------")
		//fmt.Println(s.Prices[:40])
		for i := 0; i < 40; i++ {
			fmt.Println("date is: ", s.Prices[i].T.String(), " price is: ", s.Prices[i].SP.Close)
		}

		EMA_12_period := indicator.EMA_Indicator{}
		EMA_26_period := indicator.EMA_Indicator{}
		MACD_indicator := indicator.MACD_Indicator{}

		EMA_12_period.Period = 12
		EMA_12_period.Asset_Type = s.Type
		EMA_12_period.Smoothing = 2

		EMA_26_period.Period = 26
		EMA_26_period.Asset_Type = s.Type
		EMA_26_period.Smoothing = 2

		MACD_indicator.Asset_Type = s.Type
		MACD_indicator.Smoothing_EMA = 2

		fmt.Println("----------------------------------EMA 12 indicator is below----------------------------")

		price, err := EMA_12_period.Get_Indicator(s)
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println(price[:28])
		for i := 0; i < 28; i++ {
			fmt.Println("date is: ", price[i].T.String(), " price is: ", price[i].P)
		}

		fmt.Println("----------------------------------EMA 26 indicator is below----------------------------")

		price, err = EMA_26_period.Get_Indicator(s)
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println(price[:14])
		for i := 0; i < 14; i++ {
			fmt.Println("date is: ", price[i].T.String(), " price is: ", price[i].P)
		}

		fmt.Println("----------------------------------MACD sig indicator is below----------------------------")

		sig, price, err := MACD_indicator.Get_Hist_Indicator(s)
		if err != nil {
			fmt.Println(err.Error())
		}
		//fmt.Println(sig[:5])
		for i := 0; i < 5; i++ {
			fmt.Println("date is: ", sig[i].T.String(), " price is: ", sig[i].P)
		}

		fmt.Println("----------------------------------MACD  indicator is below----------------------------")
		//fmt.Println(price[:14])
		for i := 0; i < 14; i++ {
			fmt.Println("date is: ", price[i].T.String(), " price is: ", price[i].P)
		}
	*/

	EMA_12_period := indicator.EMA_Indicator{}
	macd_indic := indicator.MACD_Indicator{}

	EMA_12_period.Asset_Type = s.Type
	EMA_12_period.Smoothing = 2
	EMA_12_period.Period = 12

	macd_indic.Asset_Type = s.Type
	macd_indic.Smoothing_EMA = 2

	test_price := s
	test_price.Prices = test_price.Prices[100:]
	EMA_12_period.Indicator_Value, _ = EMA_12_period.Calculate_Indicator(test_price)
	macd_indic.Signal_Indicator, macd_indic.MACD_Indicator_Value, err = macd_indic.Calculate_Indicator(test_price)

	fmt.Println("--------------------data----------------------")
	for i := 0; i < len(test_price.Prices); i++ {
		fmt.Println("data is: ", test_price.Prices[i].T.String(), " price is: ", test_price.Prices[i].SP)
	}

	fmt.Println("--------------------ema data----------------------")
	for i := 0; i < len(EMA_12_period.Indicator_Value); i++ {
		fmt.Println("data is: ", EMA_12_period.Indicator_Value[i].T.String(), " price is: ", EMA_12_period.Indicator_Value[i].P)
	}

	fmt.Println("---------------------macd data-----------------------------")
	for i := 0; i < len(macd_indic.MACD_Indicator_Value); i++ {
		fmt.Println("data is: ", macd_indic.MACD_Indicator_Value[i].T.String(), " price is: ", macd_indic.MACD_Indicator_Value[i].P)
	}

	fmt.Println("-----------------------macd signal data------------------------------------")
	for i := 0; i < len(macd_indic.Signal_Indicator); i++ {
		fmt.Println("data is: ", macd_indic.Signal_Indicator[i].T.String(), " price is: ", macd_indic.MACD_Indicator_Value[i].P)
	}

	//for test indicator

	c.JSON(200, s.Prices)

}
