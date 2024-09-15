package model

import "go_web/app/core"

type Article struct {
	core.ModelId
	Uuid    uint64 `gorm:"index"`
	Title   string `gorm:"size:128"`
	Content string `gorm:"type:text"`
	viewed  int64  `gorm:"default:0"`
	Status  int8   `gorm:"default:1"`
}
