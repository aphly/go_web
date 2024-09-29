package model

import (
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/core"
)

type ChatMessage struct {
	Id        core.Int64 `gorm:"primarykey" json:"id,omitempty"`
	ChatId    core.Int64 `gorm:"index" json:"chat_id,omitempty"`
	Uuid      core.Int64 `gorm:"index" json:"uuid,omitempty"`
	MsgType   int8       `json:"msg_type,omitempty"`
	ToId      core.Int64 ` json:"to_id,omitempty"`
	Message   string     `gorm:"size:1024" json:"message,omitempty"`
	Pic       string     `json:"pic,omitempty"`
	Url       string     `json:"url,omitempty"`
	Desc      string     `json:"desc,omitempty"`
	Amount    int        `json:"amount,omitempty"`
	IsDel     int64      `gorm:"default:0" json:"is_del,omitempty"`
	CreatedAt int64      `gorm:"autoUpdateTime"  json:"created_at,omitempty"`
}

func Publish(c *gin.Context, channel string, msg string) error {
	var err error
	err = app.RedisW().Publish(c, channel, msg).Err()
	return err
}

func Subscribe(c *gin.Context, channel string) (string, error) {
	sub := app.RedisW().PSubscribe(c, channel)
	msg, err := sub.ReceiveMessage(c)
	return msg.Payload, err
}

type HistoryQuery struct {
	ChatId core.Int64
	Page   int
}
