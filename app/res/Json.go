package res

import (
	"github.com/gin-gonic/gin"
)

type response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

type OptionFunc func(*response)

func Code(data int) OptionFunc {
	return func(this *response) {
		this.Code = data
	}
}

func Msg(data string) OptionFunc {
	return func(this *response) {
		this.Msg = data
	}
}

func Data(data any) OptionFunc {
	return func(this *response) {
		this.Data = data
	}
}

func Json(c *gin.Context, options ...OptionFunc) {
	res := &response{
		Code: 0,
		Msg:  "成功",
	}
	for _, v := range options {
		v(res)
	}
	c.JSON(200, res)
}
