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
	AccessTokenExpire  = 7200
	RefreshTokenExpire = 31536000
)

type User struct {
	Uuid               int64   `gorm:"primarykey"`
	Nickname           string  `gorm:"size:16"`
	AccessToken        string  `gorm:"index;size:64"`
	RefreshToken       string  `gorm:"index;size:64"`
	AccessTokenExpire  int64   `gorm:"default:0"`
	RefreshTokenExpire int64   `gorm:"default:0"`
	Avatar             *string `gorm:"size:255"`
	Remote             int8    `gorm:"default:0"`
	Status             int8    `gorm:"default:1"`
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

func (this *User) Add(uuid int64) error {
	this.Uuid = uuid
	this.Nickname = helper.RandStr(10)
	now := time.Now().Unix()
	this.AccessToken = helper.RandStr(32)
	this.RefreshToken = helper.RandStr(32, 1)
	this.AccessTokenExpire = now + 7200
	this.RefreshTokenExpire = now + 31536000
	app.DbW().Create(this)
	return nil
}

func (this *User) EnToken(token string) string {
	uuidStr := strconv.FormatInt(this.Uuid, 10)
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
	this.AccessTokenExpire = now + 30
	this.RefreshTokenExpire = now + 300
}

func (this *User) GenAccessToken(now int64) {
	this.AccessToken = helper.RandStr(32)
	this.AccessTokenExpire = now + 30
}
