package model

import "go_web/app/core"

type ChatMember struct {
	ChatId    core.Int64 `gorm:"primarykey" json:"chat_id,omitempty"`
	Uuid      core.Int64 `gorm:"primarykey" json:"uuid,omitempty"`
	IsNew     int8       `json:"is_new,omitempty"`
	UpdatedAt int64      `gorm:"index;autoCreateTime" json:"updated_at,omitempty"`
	core.Model
}
