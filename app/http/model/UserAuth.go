package model

import (
	"github.com/go-playground/validator/v10"
	"go_web/app/core"
)

type UserAuth struct {
	IdType         string `gorm:"primaryKey;type:char(16)" json:"id_type"`
	Id             string `gorm:"primaryKey;size:64" form:"id" json:"id"`
	Password       string `gorm:"size:255" form:"password"`
	Uuid           int64  `gorm:"index" json:"uuid"`
	LastIp         string `gorm:"size:64" json:"last_ip"`
	LastTime       int64  `json:"last_time"`
	Note           string `gorm:"size:255"`
	UserAgent      string `gorm:"type:text"`
	AcceptLanguage string `gorm:"size:255"`
	Verified       int8   `gorm:"default:0"`
	core.Model
}

type UserAuthForm struct {
	IdType   string
	Id       string `form:"id" binding:"required,len=11,numeric"`
	Password string `form:"password" binding:"required,min=6,max=32"`
}

func (r *UserAuthForm) GetError(err validator.ValidationErrors) string {
	str := "校验格式错误"
	for _, v := range err {
		if v.Field() == "Password" {
			switch v.Tag() {
			case "required":
				str = "请输入密码"
			case "min":
				str = "密码最小长度为6位"
			case "max":
				str = "密码最大长度为32位"
			}
		} else if v.Field() == "Id" {
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
