package test

import (
	"fmt"
	"go_web/app"
	"go_web/app/http/model"
	"testing"
)

var AutoTables = []any{
	//&model.User{},
	//&model.UserAuth{},
	//&model.Article{},
	//&model.Friend{},
	//&model.FriendRequest{},

	//&model.Chat{},
	&model.ChatMessage{},
	//&model.ChatMember{},
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
