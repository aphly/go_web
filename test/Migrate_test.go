package test

import (
	"fmt"
	"go_blog/app"
	"go_blog/app/http/model"
	"testing"
)

var AutoTables = []any{
	&model.User{},
	&model.UserAuth{},
	&model.Article{},
}

func TestMigrate(t *testing.T) {
	app.Init()
	for _, v := range AutoTables {
		err := app.DbW().AutoMigrate(v)
		if err != nil {
			fmt.Println(err)
		}
	}
	fmt.Println("AutoMigrate ok")
}
