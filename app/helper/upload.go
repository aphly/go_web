package helper

import (
	"fmt"
	"go_web/app/core"
	"strings"
)

type Upload struct {
	SizeConf  int64
	TypeConf  []any
	ConutConf int
}

func NewUpload() *Upload {
	return &Upload{
		ConutConf: 1,
		SizeConf:  1,
		TypeConf: []any{
			"image/jpeg",
			"image/png",
		},
	}
}

func (this *Upload) TypeLimit(fileTypes []string) bool {
	for _, i := range fileTypes {
		if !InArray(i, this.TypeConf) {
			return false
		}
	}
	return true
}

func (this *Upload) SizeLimit(fileSize int64) bool {
	if fileSize <= this.SizeConf*1024*1024 {
		return true
	}
	return false
}

func UuidPath(uuid core.Int64) string {
	str := fmt.Sprintf("%0*d", 20, uuid)
	splitted := make([]string, len(str)/2)
	for i := 0; i < len(str); i += 2 {
		splitted[i/2] = str[i : i+2]
	}
	return strings.Join(splitted, "/")
}
