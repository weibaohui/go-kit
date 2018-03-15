package httpkit

import (
	"strings"

	"github.com/gin-gonic/gin"
)

/**
GetServerPath 获取当前服务器的路径
*/
func GetServerPath(c *gin.Context) string {
	serverPath := strings.ToLower(strings.Split(c.Request.Proto, "/")[0]) + "://" + c.Request.Host
	return serverPath
}
