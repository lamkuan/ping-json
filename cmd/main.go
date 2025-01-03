package main

import (
	"github.com/gin-gonic/gin"
	//_ "github.com/lamkuan/ping-json/docs"
	_ "github.com/lamkuan/ping-json/docs"
	"github.com/lamkuan/ping-json/internal/controllers"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"
	"os/signal"
)

// @title Ping API 文档
// @version 1.0
// @description 一个用于执行 ping 命令的示例 API。
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @host localhost:8080
// @BasePath /
func main() {
	r := gin.Default()

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		for _ = range c {
			os.Exit(1)
		}
	}()

	// 添加 Swagger 路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := r.Group("/")
	{
		pingAPI := api.Group("/ping")
		{
			// @Summary Ping 目标地址
			// @Description 发送 ping 请求到目标 IP 地址，并返回结果。
			// @Tags Ping
			// @Param ip path string true "目标 IP 地址"
			// @Param params path string false "附加参数 (count, timeout)"
			// @Param get_latency query string false "是否获取延迟时间（yes 或 no）"
			// @Success 200 {object} map[string]interface{}
			// @Failure 400 {object} map[string]interface{}
			// @Router /ping/{ip}/{params} [get]
			pingAPI.GET("/:ip/*params", controllers.Ping)
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}
