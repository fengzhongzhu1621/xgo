package main

import (
	"archive/zip"
	"bytes"
	"errors"
	"io"
	"net/url"
	"path/filepath"
	"regexp"

	goplist "github.com/fork2fix/go-plist"
)

// 读取压缩包中的图标文件
func ParseIpaIcon(path string) (data []byte, err error) {
	iconPattern := regexp.MustCompile(`(?i)^Payload/[^/]*/icon\.png$`)
	// 打开压缩包
	r, err := zip.OpenReader(path)
	if err != nil {
		return
	}
	defer r.Close()

	// 编译压缩文件匹配图标文件
	var zfile *zip.File
	for _, file := range r.File {
		if iconPattern.MatchString(file.Name) {
			zfile = file
			break
		}
	}
	if zfile == nil {
		err = errors.New("icon.png file not found")
		return
	}

	// 打开图标文件
	plreader, err := zfile.Open()
	if err != nil {
		return
	}
	defer plreader.Close()
	// 用于读取 io.Reader 接口的所有数据到一个字节切片（[]byte）中。
	// 如果读取成功，函数会返回读取到的数据和 nil 错误；如果读取失败，则会返回一个非空的错误。
	return io.ReadAll(plreader)
}

// 读取压缩包中的 Info.plist 文件
func ParseIPA(path string) (plinfo *plistBundle, err error) {
	plistre := regexp.MustCompile(`^Payload/[^/]*/Info\.plist$`)
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var plfile *zip.File
	for _, file := range r.File {
		if plistre.MatchString(file.Name) {
			plfile = file
			break
		}
	}
	if plfile == nil {
		err = errors.New("the Info.plist file not found")
		return nil, err
	}
	plreader, err := plfile.Open()
	if err != nil {
		return nil, err
	}
	defer plreader.Close()

	// 读取 Info.plist 文件
	// plist 是一种用于存储用户设置和应用程序数据的文件格式，主要在 macOS 和 iOS 系统上使用。
	buf := make([]byte, plfile.FileInfo().Size())
	_, err = io.ReadFull(plreader, buf)
	if err != nil {
		return nil, err
	}

	// 解析 Info.plist 文件内容到结构体中
	// 用于创建一个新的 io.Reader 接口，该接口从给定的字节切片（[]byte）中读取数据。
	// 这使得你可以将字节切片作为文件或网络连接等具有 io.Reader 接口的对象来读取。
	dec := goplist.NewDecoder(bytes.NewReader(buf))
	plinfo = new(plistBundle)
	err = dec.Decode(plinfo)
	return plinfo, err
}

type plistBundle struct {
	CFBundleIdentifier  string `plist:"CFBundleIdentifier"`
	CFBundleVersion     string `plist:"CFBundleVersion"`
	CFBundleDisplayName string `plist:"CFBundleDisplayName"`
	CFBundleName        string `plist:"CFBundleName"`
	CFBundleIconFile    string `plist:"CFBundleIconFile"`
	CFBundleIcons       struct {
		CFBundlePrimaryIcon struct {
			CFBundleIconFiles []string `plist:"CFBundleIconFiles"`
		} `plist:"CFBundlePrimaryIcon"`
	} `plist:"CFBundleIcons"`
}

// ref: https://gist.github.com/frischmilch/b15d81eabb67925642bd#file_manifest.plist
type plAsset struct {
	Kind string `plist:"kind"`
	URL  string `plist:"url"`
}

type plItem struct {
	Assets   []*plAsset `plist:"assets"`
	Metadata struct {
		BundleIdentifier string `plist:"bundle-identifier"`
		BundleVersion    string `plist:"bundle-version"`
		Kind             string `plist:"kind"`
		Title            string `plist:"title"`
	} `plist:"metadata"`
}

type downloadPlist struct {
	Items []*plItem `plist:"items"`
}

func GenerateDownloadPlist(baseURL *url.URL, ipaPath string, plinfo *plistBundle) ([]byte, error) {
	dp := new(downloadPlist)
	item := new(plItem)
	baseURL.Path = ipaPath
	ipaUrl := baseURL.String()
	item.Assets = append(item.Assets, &plAsset{
		Kind: "software-package",
		URL:  ipaUrl,
	})

	iconFiles := plinfo.CFBundleIcons.CFBundlePrimaryIcon.CFBundleIconFiles
	if iconFiles != nil && len(iconFiles) > 0 {
		baseURL.Path = "/-/unzip/" + ipaPath + "/-/**/" + iconFiles[0] + ".png"
		imgUrl := baseURL.String()
		item.Assets = append(item.Assets, &plAsset{
			Kind: "display-image",
			URL:  imgUrl,
		})
	}

	item.Metadata.Kind = "software"

	item.Metadata.BundleIdentifier = plinfo.CFBundleIdentifier
	item.Metadata.BundleVersion = plinfo.CFBundleVersion
	item.Metadata.Title = plinfo.CFBundleName
	if item.Metadata.Title == "" {
		item.Metadata.Title = filepath.Base(ipaUrl)
	}

	dp.Items = append(dp.Items, item)
	// 将数据结构编码为 plist 格式，并且以缩进的方式输出。
	data, err := goplist.MarshalIndent(dp, goplist.XMLFormat, "    ")
	return data, err
}
