package main

import (
	"algotrading/backtest"
	"algotrading/db"
	"algotrading/handler"
	"algotrading/middlewares"

	"github.com/gin-gonic/gin"
)

func main() {

	//for access-control-allow
	router := gin.New()
	router.Use(middlewares.Cors())
	//for restful api
	api_router := router.Group("/api")
	//asset api
	stock_api := api_router.Group("/stock")
	stock_api.POST("/query", handler.Asset_Query)

	d := db.Db_main()
	defer d.Close()
	backtest.Backtest_Main(d)

	router.Run("192.168.152.215:8080")

}
