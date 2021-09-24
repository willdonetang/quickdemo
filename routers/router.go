package routers

import (
	"bytes"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	_ "quickdemo/docs"
	"quickdemo/pkg/app"
	"quickdemo/pkg/logf"
	"quickdemo/pkg/util/const"
	"strings"
	"time"
)

// InitRouter initialize routing information
func InitRouter() *gin.Engine {
	r := gin.New()
	//日志中间件,所有的异常捕获,logrus
	r.Use(gin.Recovery(), initLog(), cors())

	r.GET("/test", test)

	return r
}

func test(c *gin.Context) {
	app.SuccessResp(c, "test")
	//c.String(200, "test")
}

//跨域中间件
func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, "+utilconst.HeaderToken)
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type, "+utilconst.HeaderToken)
			c.Header("Access-Control-Allow-Credentials", "false")
			c.Set("content-type", "application/json")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func initLog() gin.HandlerFunc {

	return func(c *gin.Context) {

		// 开始时间
		startTime := time.Now()
		// 请求路由
		path := c.Request.RequestURI

		// 排除文件上传的请求体打印
		isFormData := strings.Contains(c.Request.Header.Get("Content-Type"), "multipart/form-data")
		// requestBody
		var requestBody []byte
		if !isFormData {
			requestBody, _ = c.GetRawData()
			c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))
		}

		//处理请求
		c.Next()
		// 处理结果
		result, exists := c.Get(utilconst.LogResponse)
		if exists {
			result = result.(*app.Response)
		}

		// 执行时间
		latencyTime := time.Since(startTime)
		// 请求方式
		reqMethod := c.Request.Method
		// http状态码
		statusCode := c.Writer.Status()
		// 请求IP
		clientIP := c.ClientIP()
		//token := c.GetHeader(tool.HeaderToken)
		// 日志格式
		logf.InfoWithFields(logrus.Fields{
			"req_body":     string(requestBody),
			"http_code":    statusCode,
			"latency_time": fmt.Sprintf("%13v", latencyTime),
			"ip":           clientIP,
			"method":       reqMethod,
			"path":         path,
			"result":       result,
		})
	}
}
