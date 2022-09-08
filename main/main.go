package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_dev/gogin/jwt"
	"go_dev/gogin/middleware"
	"go_dev/gogin/utils"
	"io/ioutil"
	"net/http"
)

func CorsMiddleware() gin.HandlerFunc {
	return func(gctx *gin.Context) {
		method := gctx.Request.Method
		origin := gctx.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			// 接收客户端发送的origin
			gctx.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			// 服务器支持的所有跨域请求的方法
			gctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			// 允许跨域设置可以返回其他子段，可以自定义字段
			gctx.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, Content-Type, X-CSRF-Token, Token, session")
			// 允许浏览器（客户端）可以解析的头部
			gctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			// 设置缓存时间
			gctx.Header("Access-Control-Max-Age", "172800")
			// 允许客户端传递校验信息比如 cookie
			gctx.Header("Access-Control-Allow-Credentials", "true")
		}

		// 允许类型校验，不经过网关，放行所有 OPTIONS 请求
		if method == "OPTIONS" {
			gctx.JSON(http.StatusOK, "ok!")
			gctx.Abort()
			return
		}

		gctx.Next()
	}
}

type User struct {
	UserName  string `json:"username"`
	Password  string `json:"password"`
	EmailCode string `json:"emailcode"`
}

// 解析token并在token中获取用户信息
type Permissions struct {
	Id      int    `json:"id"`
	Menu    string `json:"menu"`
	Url     string `json:"url"`
	Methods string `json:"method"`
}

type UserInfo struct {
	UserName        string        `json:"username"`
	Password        string        `json:"password"`
	PermissionsList []Permissions `json:"permissions"`
}

var (
	userinfo    UserInfo
	permissions []Permissions
)

func GetUserInfo(context *gin.Context) {
	// 判断header是否存在admin-token
	token := context.Request.Header.Get("token")
	fmt.Println("token", token)
	// 判断token是否存在
	if token == "" {
		context.JSON(501, gin.H{
			"message": "token不存在",
		})
	}
	permissions = append(permissions, Permissions{Id: 1, Menu: "用户管理", Url: "/admin/user", Methods: "GET,POST,DELETE,PUT"})
	permissions = append(permissions, Permissions{Id: 2, Menu: "角色管理", Url: "/admin/role", Methods: "GET,POST,DELETE,PUT"})
	permissions = append(permissions, Permissions{Id: 3, Menu: "权限管理", Url: "/admin/permission", Methods: "GET,POST,DELETE,PUT"})
	permissions = append(permissions, Permissions{Id: 4, Menu: "系统管理", Url: "/admin/system", Methods: "GET,POST,DELETE,PUT"})
	userinfo = UserInfo{UserName: "admin", Password: "admin", PermissionsList: permissions}
	context.JSON(200, gin.H{
		"permissions": userinfo,
	})
}

func GetEmailCode(c *gin.Context) {
	data, _ := ioutil.ReadAll(c.Request.Body)
	var user User
	_ = json.Unmarshal(data, &user)
	// 模拟通过用户名和密码从数据库查询该用户的邮箱
	if user.UserName != "admin" && user.Password != "admin" {
		c.JSON(http.StatusOK, gin.H{
			"message": "邮件发送失败",
		})
		return
	}
	mail := []string{"573937686@qq.com"}
	// 随机生成4位验证码
	code := utils.GetRandomCode(4)
	err := utils.SendMail(mail, "系统验证码", fmt.Sprintf("验证码：%s", code))
	// 将code保存至cookie中
	session := sessions.Default(c)
	session.Set("code", code)
	_ = session.Save()
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"message": "邮件发送失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "邮件发送成功",
	})
}

func main() {

	// 初始化数据库
	//common.DBInit()

	r := gin.Default()
	r.SetTrustedProxies([]string{})
	r.Use(middleware.Session("mark"))
	// 使用中间件
	//r.Use(Cors()) // 允许跨域
	r.POST("/admin/login", func(context *gin.Context) {
		fmt.Println(context.Request.Header.Get("X-FORWARD-FOR"))
		data, _ := ioutil.ReadAll(context.Request.Body)
		fmt.Println(string(data))
		// 序列化data
		var user User
		_ = json.Unmarshal(data, &user)
		// 获取用户信息并生成token
		ss, err := jwt.CreateToken(user.UserName, user.Password)
		if err != nil {
			context.JSON(http.StatusOK, gin.H{
				"code":    http.StatusInternalServerError,
				"message": "登录失败",
			})
			return
		}
		// 从session获取code
		session := sessions.Default(context)
		code := session.Get("code")
		session.Delete("code")
		_ = session.Save()
		fmt.Println("code", code)

		// 判断验证码
		//if code != user.EmailCode {
		//	context.JSON(401, gin.H{
		//		"code":    401,
		//		"message": "验证码错误",
		//	})
		//	return
		//}

		if user.UserName == "admin" && user.Password == "admin" {
			context.JSON(200, gin.H{
				"message": "用户登入成功",
				"code":    200,
				"token":   ss,
			})
			return
		} else {
			context.JSON(401, gin.H{
				"message": "用户登入失败",
				"code":    401,
			})
			return
		}
	})

	r.POST("/admin/userInfo", GetUserInfo)
	r.POST("/admin/sendemailcode", GetEmailCode)

	apiGroup := r.Group("/api").Use(middleware.JwtMiddleware())
	{
		apiGroup.GET("/about", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "about",
			})
		})
	}

	err := r.Run(":8080")
	if err != nil {
		fmt.Println(err)
	}
}
