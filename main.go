package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"quickdemo/crontab"
	"quickdemo/models"
	"quickdemo/pkg/gredis"
	"quickdemo/pkg/logf"
	"quickdemo/pkg/setting"
	"quickdemo/routers"
	vali "quickdemo/validation"
)

func init() {
	setting.Setup()
	logf.Setup()
	models.Setup()
	gredis.Setup()
	crontab.Setup()
	vali.InitValidation()
}

// @title Golang Gin API
// @version 1.0
// @description An example of gin
// @termsOfService https://github.com/EDDYCJY/go-gin-example
// @license.name MIT
// @license.url https://Pay/blob/master/LICENSE
func main() {
	gin.SetMode(setting.ServerSetting.RunMode)
	r := routers.InitRouter()
	endPoint := fmt.Sprintf(":%d", setting.ServerSetting.HttpPort)
	maxHeaderBytes := 1 << 20

	server := &http.Server{
		Addr:           endPoint,
		Handler:        r,
		MaxHeaderBytes: maxHeaderBytes,
	}

	log.Printf("[info] start http server listening %s", endPoint)

	server.ListenAndServe()
}
