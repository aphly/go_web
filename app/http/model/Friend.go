package model

import "go_web/app/core"

type Friend struct {
	Uuid       core.Int64 `gorm:"primarykey" json:"uuid,omitempty"`
	FriendUuid core.Int64 `gorm:"primarykey" json:"friend_uuid,omitempty"`
	Status     int8       `json:"status,omitempty"`
	core.Model
	FriendUser User `gorm:"foreignKey:FriendUuid"`
}
