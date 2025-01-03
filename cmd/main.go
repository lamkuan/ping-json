package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"regexp"
	"strconv"
	"strings"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/lamkuan/ping-json/internal/ping"
	
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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
			pingAPI.GET("/:ip/*params", func(c *gin.Context) {
				var err error

				ip := c.Param("ip")
				params := c.Param("params")
				get_latency := c.Query("get_latency")

				parts := strings.Split(strings.TrimPrefix(params, "/"), "/")
				var countStr, timeoutStr string

				if len(parts) > 0 {
					countStr = parts[0]
				}
				if len(parts) > 1 {
					timeoutStr = parts[1]
				}

				count := 5
				if countStr != "" {
					count, err = strconv.Atoi(countStr)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
						return
					}
				}

				timeout := 3600 * time.Second
				if timeoutStr != "" {
					x, err := strconv.Atoi(timeoutStr)
					timeout = time.Duration(x)
					if err != nil {
						c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid timeout parameter"})
						return
					}
				}

				result, _ := ping.Ping(map[string]interface{}{"ip": ip, "count": count, "timeout": timeout})
				var latencyList []string

				if get_latency == "yes" {
					re := regexp.MustCompile(`time=(\d+\.?\d*)`)
					matches := re.FindAllStringSubmatch(result, -1)
					latencyList = make([]string, 0, 1000)
					if len(matches) > 0 {
						for _, match := range matches {
							latencyList = append(latencyList, fmt.Sprintf("%s\n", match[1]))
						}
					}

					c.JSON(http.StatusOK, gin.H{
						"result": result,
						"status": http.StatusOK,
						"times":  latencyList,
					})
					return
				}

				c.JSON(http.StatusOK, gin.H{
					"result": result,
					"status": http.StatusOK,
				})
			})
		}
	}

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err.Error())
	}
}
