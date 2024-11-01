package service

import (
	"go_web/app"
	"os"
)

func UploadPath(path string, remote int8) string {
	if path == "" {
		return ""
	}
	if remote == 1 {
		if app.Config.Oss.IsCname == "true" {
			return app.Config.Oss.Url + path
		} else {
			return app.Config.Oss.Endpoint + path
		}
	} else if remote == 0 {
		return app.Config.Http.Host + path
	} else {
		return ""
	}
}

func UploadDel(path string, remote int8) error {
	if remote == 1 {
		if app.Config.Oss.IsCname == "true" {
			return nil
		} else {
			return nil
		}
	} else if remote == 0 {
		if err := os.Remove("." + path); err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}
