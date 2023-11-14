package handler

import (
	"algotrading/asset"
	"algotrading/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Asset_Query(c *gin.Context) {
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

	err = s.Get_Price()
	if err != nil {
		c.JSON(500, gin.H{
			"code":    500,
			"message": "Server Internal Error",
		})
		logger.Info.Fatal(err)
		return
	}

	//fmt.Println(s.Prices)
	c.JSON(200, s.Prices)

}
