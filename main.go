package main

import (
	"go_web/app"
	"go_web/app/http/router"
)

func main() {
	app.Init()
	server := router.Reg()
	err := server.Run(app.Config.Http.Listen)
	if err != nil {
		panic(err)
	}
}
