package my

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"go_web/app"
	"go_web/app/helper"
	"go_web/app/http/model"
	"go_web/app/http/service"
	"go_web/app/res"
	"os"
	"strconv"
	"strings"
	"time"
)

func Avatar(c *gin.Context) {
	file, _ := c.FormFile("file")
	stringArr := file.Header["Content-Type"]
	var upload = helper.NewUpload()
	if !upload.SizeLimit(file.Size) {
		res.Json(c, res.Code(1), res.Msg("大小超过"+strconv.FormatInt(upload.SizeConf, 10)+"M"))
		return
	}

	if !upload.TypeLimit(stringArr) {
		res.Json(c, res.Code(2), res.Msg("格式只支持"+helper.SliceJoin(upload.TypeConf, "、")))
		return
	}
	var now = time.Now()
	var timeDir = "/public/upload/avatar/" + fmt.Sprintf("%d%d/%d/%d%d/", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute())
	err := os.MkdirAll("."+timeDir, 0755)
	if err != nil {
		res.Json(c, res.Code(3), res.Msg("文件夹错误"))
		return
	}
	arr := strings.Split(file.Filename, ".")
	if len(arr) != 2 {
		res.Json(c, res.Code(4), res.Msg("上传文件格式错误"))
		return
	}
	var savePath = timeDir + helper.RandStr(16) + "." + arr[1]
	var uploadPath = "." + savePath
	err = c.SaveUploadedFile(file, uploadPath)
	if err != nil {
		res.Json(c, res.Code(5), res.Msg("上传错误"))
		return
	}
	getUser, _ := c.Get("user")
	me, _ := getUser.(model.User)

	var user model.User
	err = app.DbW().Model(&user).Select("avatar", "remote").Where("uuid = ?", me.Uuid).
		Updates(model.User{Avatar: savePath, Remote: 0}).Error
	if err != nil {
		res.Json(c, res.Code(6), res.Msg("保存错误"))
		return
	}
	service.UploadDel(me.Avatar, me.Remote)
	res.Json(c, res.Data(gin.H{
		"avatar": service.UploadPath(user.Avatar, user.Remote),
	}))
	return
}

func Pics(c *gin.Context) {
	file, _ := c.FormFile("file")
	dst := "./" + file.Filename
	fmt.Println(dst)
	res.Json(c, res.Data(gin.H{
		"avatar": "xxx",
	}))
	return
}
