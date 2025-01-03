package controllers

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lamkuan/ping-json/internal/ping"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Ping handles ICMP ping requests.
//
// @Summary      ICMP Ping
// @Description  Send an ICMP ping request to a specified IP address and retrieve results, including latency if requested.
// @Tags         Ping
// @Param        ip             path      string  true   "Target IP Address"
// @Param        params         path      string  false  "Ping parameters in the format: {count}/{timeout}. Default count is 5, and default timeout is 3600 seconds."
// @Param        get_latency    query     string  false  "Specify 'yes' to include latency data in the response."
// @Success      200            {object}  map[string]interface{}   "Response with ping result and latency data (if requested)"
// @Failure      400            {object}  map[string]interface{}   "Invalid parameters"
// @Router       /ping/{ip}/{params} [get]
func Ping(c *gin.Context) {
	var err error

	ip := c.Param("ip")
	params := c.Param("params")
	get_latency := c.Query("get_latency")

	parts := strings.Split(strings.TrimPrefix(params, "/"), "/")
	var countStr, timeoutStr string

	fmt.Println(countStr)

	if len(parts) > 0 {
		countStr = parts[0]
	}
	if len(parts) > 1 {
		timeoutStr = parts[1]
	}

	fmt.Println(countStr, timeoutStr, get_latency)

	count := 5
	if countStr != "" && countStr != "undefined" {
		count, err = strconv.Atoi(countStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid count parameter"})
			return
		}
	}

	timeout := 3600 * time.Second
	if timeoutStr != "" && timeoutStr != "undefined" {
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
}
