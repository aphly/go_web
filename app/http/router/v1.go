package router

import (
	"github.com/gin-gonic/gin"
	"go_web/app/http/controller/user"
	"go_web/app/http/middleware"
)

func v1(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/register", user.Register)
		v1.POST("/login", user.Login)
		v1.GET("/refresh_token", user.RefreshToken)
		authorized := v1.Group("")
		authorized.Use(middleware.AuthHandler())
		{
			authorized.GET("/t1", user.Test)
		}
		ipLimit := v1.Group("")
		ipLimit.Use(middleware.IpLimitHandler())
		{
			ipLimit.GET("/t2", user.Test1)
		}
	}
}
