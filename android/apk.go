package android

import (
	"log"

	"github.com/shogo82148/androidbinary/apk"
)

type ApkInfo struct {
	PackageName  string `json:"packageName"`
	MainActivity string `json:"mainActivity"`
	Version      struct {
		Code int32  `json:"code"`
		Name string `json:"name"`
	} `json:"version"`
}

// ParseApkInfo path should be absolute
func ParseApkInfo(path string) (ai *ApkInfo) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("parse-apk-info panic:", err)
		}
	}()
	apkf, err := apk.OpenFile(path)
	if err != nil {
		return nil
	}
	ai = &ApkInfo{}
	ai.MainActivity, _ = apkf.MainActivity()
	ai.PackageName = apkf.PackageName()
	code, _ := apkf.Manifest().VersionCode.Int32()
	ai.Version.Code = code
	name, _ := apkf.Manifest().VersionName.String()
	ai.Version.Name = name

	return ai
}
