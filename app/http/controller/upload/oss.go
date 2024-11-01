package upload

import (
	"github.com/gin-gonic/gin"
	"go_web/app/http/service"
	"go_web/app/res"
)

func PolicyToken(c *gin.Context) {
	policyToken, err := service.OssGetPolicyToken()
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(err.Error()))
		return
	}
	res.Json(c, res.Data(policyToken))
	return
}

func OssPost(c *gin.Context) {
	bytePublicKey, err := service.GetPublicKey(c)
	if err != nil {
		res.Json(c, res.Code(1), res.Msg(err.Error()))
		return
	}
	byteAuthorization, err := service.GetAuthorization(c)
	if err != nil {
		res.Json(c, res.Code(2), res.Msg(err.Error()))
		return
	}

	byteMD5, err := service.GetMD5FromNewAuthString(c)
	if err != nil {
		res.Json(c, res.Code(3), res.Msg(err.Error()))
		return
	}

	if service.OssVerifySignature(bytePublicKey, byteMD5, byteAuthorization) {
		res.Json(c)
	} else {
		res.Json(c, res.Code(4), res.Msg("400"))
	}
}
