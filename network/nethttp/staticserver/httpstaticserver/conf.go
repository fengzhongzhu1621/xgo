package httpstaticserver

import (
	"log"
	"os"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/crypto/compress/zipfile"
	"github.com/fengzhongzhu1621/xgo/file"
	"github.com/fengzhongzhu1621/xgo/network/nethttp/staticserver"
	"gopkg.in/yaml.v2"
)

// 构造默认配置对象
func (s *HTTPStaticServer) defaultAccessConf() staticserver.AccessConf {
	return staticserver.AccessConf{
		Upload: s.Upload,
		Delete: s.Delete,
	}
}

// readAccessConf 从路径所在的目录和祖先目录读取静态文件关联的配置文件，转换为结构体
func (s *HTTPStaticServer) readAccessConf(realPath string) (ac staticserver.AccessConf) {
	// 获得静态文件的相对路径，先判断静态文件是否在根目录
	relativePath, err := filepath.Rel(s.Root, realPath)
	if err != nil || relativePath == "." || relativePath == "" { // actually relativePath is always "." if root == realPath
		// 构造默认配置对象
		ac = s.defaultAccessConf()
		// 如果静态文件在根目录，则读取根目录的配置文件
		realPath = s.Root
	} else {
		// 如果静态文件不在根目录，
		parentPath := filepath.Dir(realPath)
		// 向上遍历获取祖先目录下的配置文件，先将根据目录下的配置转换为结构体，然后在将子目录下的配置文件转换为结构体并覆盖父配置
		ac = s.readAccessConf(parentPath)
	}
	if file.IsFile(realPath) {
		realPath = filepath.Dir(realPath)
	}

	// 获得配置文件的路径
	cfgFile := filepath.Join(realPath, zipfile.YAMLCONF)
	// 读取配置文件
	data, err := os.ReadFile(cfgFile)
	if err != nil {
		if os.IsNotExist(err) {
			return ac
		}
		log.Printf("Err read .ghs.yml: %v", err)
	}
	// 将yaml文件转换为结构体对象
	err = yaml.Unmarshal(data, &ac)
	if err != nil {
		log.Printf("Err format .ghs.yml: %v", err)
	}

	return ac
}
