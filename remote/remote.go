// Copyright © 2015 Steve Francia <spf@spf13.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package remote integrates the remote features of Viper.
package remote

import (
	"bytes"
	"io"
	"os"

	crypt "github.com/sagikazarmark/crypt/config"
)

var _ RemoteConfigFactory = (*RemoteConfigProvider)(nil)

type RemoteConfigProvider struct{}

// Get 根据路径获取配置值.
func (rc RemoteConfigProvider) Get(rp RemoteProvider) (io.Reader, error) {
	// 创建远程key/value存储管理器客户端
	cm, err := GetConfigManager(rp)
	if err != nil {
		return nil, err
	}
	// 根据路径获得配置值
	b, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(b), nil
}

// Watch 监听指定路径，返回对应的值，实际上是通过Get实现的.
func (rc RemoteConfigProvider) Watch(rp RemoteProvider) (io.Reader, error) {
	cm, err := GetConfigManager(rp)
	if err != nil {
		return nil, err
	}
	// 根据路径获得配置值
	resp, err := cm.Get(rp.Path())
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(resp), nil
}

// WatchChannel 监听指定路径，返回管道.
func (rc RemoteConfigProvider) WatchChannel(rp RemoteProvider) (<-chan *RemoteResponse, chan bool) {
	// 创建 key/value 存储客户端
	cm, err := GetConfigManager(rp)
	if err != nil {
		return nil, nil
	}
	// 给quitwc 发送 true，停止 cm.Watch 中的 goroute
	quitwc := make(chan bool)
	quit := make(chan bool)
	// 接受监听收到的结果
	responsCh := make(chan *RemoteResponse)
	// 监听路径，返回响应管道
	cryptoResponseCh := cm.Watch(rp.Path(), quit)
	// need this function to convert the Channel response form crypt.Response to viper.Response
	go func(cr <-chan *crypt.Response, vr chan<- *RemoteResponse, quitwc <-chan bool, quit chan<- bool) {
		for {
			select {
			case <-quitwc:
				// 停止 cm.Watch
				quit <- true
				return
			case resp := <-cr:
				// 获得监听数据
				vr <- &RemoteResponse{
					Error: resp.Error,
					Value: resp.Value,
				}
			}
		}
	}(cryptoResponseCh, responsCh, quitwc, quit)

	return responsCh, quitwc
}

// 创建远程key/value存储管理器客户端.
func GetConfigManager(rp RemoteProvider) (crypt.ConfigManager, error) {
	// 声明一个远程配置管理器
	var cm crypt.ConfigManager
	var err error

	// 获得 secretkeyring
	// secretkeyring is the filepath to your openpgp secret keyring.  e.g. /etc/secrets/myring.gpg
	if rp.SecretKeyring() != "" {
		// 读取密钥文件内容
		var kr *os.File
		kr, err = os.Open(rp.SecretKeyring())
		if err != nil {
			return nil, err
		}
		defer kr.Close()

		// 创建远程配置管理器客户端，使用公钥加密
		switch rp.Provider() {
		case "etcd":
			cm, err = crypt.NewEtcdConfigManager([]string{rp.Endpoint()}, kr)
		case "firestore":
			cm, err = crypt.NewFirestoreConfigManager([]string{rp.Endpoint()}, kr)
		default:
			cm, err = crypt.NewConsulConfigManager([]string{rp.Endpoint()}, kr)
		}
	} else {
		switch rp.Provider() {
		case "etcd":
			cm, err = crypt.NewStandardEtcdConfigManager([]string{rp.Endpoint()})
		case "firestore":
			cm, err = crypt.NewStandardFirestoreConfigManager([]string{rp.Endpoint()})
		default:
			cm, err = crypt.NewStandardConsulConfigManager([]string{rp.Endpoint()})
		}
	}
	if err != nil {
		return nil, err
	}
	return cm, nil
}
