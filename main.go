package main

import (
	"algotrading/global"
	"algotrading/handler"
	"algotrading/logger"
	"algotrading/middlewares"

	//"net/http"
	"syscall"

	//"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"golang.org/x/term"
)

func main() {

	//for api key
	fmt.Println("Please Input alphavantage private key")
	//input key
	key, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		logger.Error.Fatal(err.Error())
	}
	//convert []byte to string
	global.Api_Key = string(key)
	//check if key is null
	if global.Api_Key == "" {
		logger.Error.Fatal("Input error")
	}

	//for access-control-allow
	router := gin.New()
	router.Use(middlewares.Cors())
	//for restful api
	api_router := router.Group("/api")
	//asset api
	stock_api := api_router.Group("/stock")
	stock_api.POST("/query", handler.Asset_Query)

	router.Run("192.168.152.215:8080")

}
