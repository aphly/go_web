package model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/core"
	"go_web/app/core/crypt"
	"go_web/app/helper"
	"strconv"
	"time"
)

const (
	AccessTokenExpire  = 300
	RefreshTokenExpire = 31536000
)

type User struct {
	Uuid               core.Int64 `gorm:"primarykey" json:"uuid,omitempty"`
	Nickname           string     `gorm:"size:16" json:"nickname,omitempty"`
	AccessToken        string     `gorm:"index;size:64" json:"access_token,omitempty"`
	RefreshToken       string     `gorm:"index;size:64" json:"refresh_token,omitempty"`
	AccessTokenExpire  int64      `gorm:"default:0" json:"-"`
	RefreshTokenExpire int64      `gorm:"default:0" json:"-"`
	Avatar             *string    `gorm:"size:255" json:"avatar,omitempty"`
	Remote             int8       `gorm:"default:0" json:"remote,omitempty"`
	Status             int8       `gorm:"default:1" json:"status,omitempty"`
	core.Model
}

func (this User) GetToken(c *gin.Context) (error, string) {
	token := c.Request.Header.Get("Authorization")
	if token == "" {
		return errors.New("No token"), ""
	}
	if len(token) < 7 || token[:6] != "Bearer " {
		return errors.New("Invalid token"), ""
	}
	return nil, token[7:]
}

func (this *User) Add(uuid core.Int64) error {
	this.Uuid = uuid
	this.Nickname = helper.RandStr(10)
	now := time.Now().Unix()
	this.AccessToken = helper.RandStr(32)
	this.RefreshToken = helper.RandStr(32, 1)
	this.AccessTokenExpire = now + AccessTokenExpire
	this.RefreshTokenExpire = now + RefreshTokenExpire
	app.DbW().Create(this)
	return nil
}

func (this *User) EnToken(token string) string {
	uuidStr := strconv.FormatInt(int64(this.Uuid), 10)
	en, _ := crypt.AesEn(uuidStr + "_" + token)
	return en
}

func (this *User) DeToken(token string) (string, error) {
	de, err := crypt.AesDe(token)
	if err != nil {
		return "", err
	}
	return de, nil
}

func (this *User) GenToken() {
	now := time.Now().Unix()
	this.AccessToken = helper.RandStr(32)
	this.RefreshToken = helper.RandStr(32, 1)
	this.AccessTokenExpire = now + AccessTokenExpire
	this.RefreshTokenExpire = now + RefreshTokenExpire
}

func (this *User) GenAccessToken(now int64) {
	this.AccessToken = helper.RandStr(32)
	this.AccessTokenExpire = now + AccessTokenExpire
}
