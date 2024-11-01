package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_web/app/http/middleware"
)

func Reg() *gin.Engine {
	//f, _ := os.Create("gin.log")
	//gin.DefaultWriter = io.MultiWriter(f)
	fmt.Println("gin")
	r := gin.New()
	r.Use(middleware.CrosHandler())
	r.Use(middleware.PanicHandler())
	r.Use(middleware.LogHandler())
	v1(r)
	return r
}
