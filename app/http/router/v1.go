package router

import (
	"github.com/gin-gonic/gin"
	"go_web/app/http/controller/article"
	"go_web/app/http/controller/chat"
	"go_web/app/http/controller/friend"
	"go_web/app/http/controller/user"
	"go_web/app/http/middleware"
)

func v1(router *gin.Engine) {
	v1 := router.Group("/v1")
	{
		v1.POST("/register", user.Register)
		v1.POST("/login", user.Login)
		v1.GET("/refresh_token", user.RefreshToken)
		//v1.GET("/chat_message", chat.Message)

		authorized := v1.Group("")
		authorized.Use(middleware.AuthHandler())
		{
			authorized.GET("/friend", friend.List)
			authorized.POST("/friend_request_add", friend.RequestAdd)
			authorized.GET("/friend_request", friend.RequestList)
			authorized.POST("/friend_request_op", friend.RequestOp)

			authorized.GET("/chat_list", chat.List)
			authorized.GET("/chat_message", chat.Message)
		}
		ipLimit := v1.Group("")
		ipLimit.Use(middleware.IpLimitHandler())
		{
			ipLimit.GET("/t2", article.List)
		}
	}
}
