package model

import (
	"github.com/go-playground/validator/v10"
	"go_web/app/core"
)

type FriendRequest struct {
	core.ModelId
	FromUuid core.Int64 `gorm:"index" json:"from_uuid,omitempty"`
	FromName string     `json:"from_name,omitempty"`
	ToUuid   core.Int64 `gorm:"index" json:"to_uuid,omitempty"`
	ToName   string     `json:"to_name,omitempty"`
	Status   int8       `json:"status,omitempty"`
	Note     string     `json:"note,omitempty"`
}

type FriendRequestForm struct {
	Id string `form:"id" binding:"required,len=11,numeric"`
}

func (this *FriendRequestForm) GetError(err validator.ValidationErrors) string {
	str := "校验格式错误"
	for _, v := range err {
		if v.Field() == "Id" {
			switch v.Tag() {
			case "required":
				str = "请输入手机号码"
			case "len":
				str = "手机号码必须11位"
			case "numeric":
				str = "手机号码必须数字"
			}
		} else {
			return v.Field() + " " + v.Tag() + "格式错误"
		}
	}
	return str
}

type FriendRequestOpForm struct {
	Id     uint `form:"id" binding:"required,numeric"`
	Status int8 `form:"status" binding:"required,numeric"`
}

func (this *FriendRequestOpForm) GetError(err validator.ValidationErrors) string {
	str := "校验格式错误"
	for _, v := range err {
		if v.Field() == "Id" {
			switch v.Tag() {
			case "required":
				str = "参数错误"
			case "numeric":
				str = "参数错误"
			}
		} else {
			return "参数错误"
		}
	}
	return str
}
