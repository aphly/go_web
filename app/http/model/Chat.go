package model

import "go_web/app/core"

type Chat struct {
	Id          core.Int64 `gorm:"primarykey" json:"id,omitempty"`
	Uuid        core.Int64 `gorm:"index"  json:"uuid,omitempty"`
	Subject     string     `gorm:"size=16"  json:"subject,omitempty"`
	MinMax      string     `gorm:"index;size:41"  json:"min_max,omitempty"`
	LastMessage string     `gorm:"size:1024"  json:"last_message,omitempty"`
	UpdatedAt   int64      `gorm:"index;autoCreateTime" json:"updated_at,omitempty"`
	CreatedAt   int64      `gorm:"autoUpdateTime" json:"created_at,omitempty"`
}
