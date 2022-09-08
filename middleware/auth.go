package middleware

import (
	"github.com/gin-gonic/gin"
	"go_dev/gogin/jwt"
	"net/http"
	"time"
)

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		//从请求头中获取token
		tokenStr := c.Request.Header.Get("token")
		// token不存在
		if tokenStr == "" {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "token不存在"})
			c.Abort() //阻止执行
			return
		}
		//token格式错误
		//tokenSlice := strings.SplitN(tokenStr, ".", 2)
		//if len(tokenSlice) != 2 {
		//	c.JSON(http.StatusOK, gin.H{"code": 0, "message": "token格式错误"})
		//	c.Abort() //阻止执行
		//	return
		//}
		//验证token
		tokenStruck, ok := jwt.CheckToken(tokenStr)
		if !ok {
			c.JSON(http.StatusOK, gin.H{"code": 100, "message": "token不正确"})
			c.Abort() //阻止执行
			return
		}
		//token超时
		if time.Now().Unix() > tokenStruck.ExpiresAt {
			c.JSON(http.StatusOK, gin.H{"code": 0, "message": "token已过期,请重新登录"})
			c.Abort() //阻止执行
			return
		}
		c.Set("username", tokenStruck.Username)
		c.Set("password", tokenStruck.Password)
		c.Next()
	}
}
