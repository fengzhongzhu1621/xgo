// Copyright 2014 Manu Martinez-Almeida.  All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package render

import "net/http"

// Render interface is to be implemented by JSON, XML, HTML, YAML and so on.
type Render interface {
	// Render writes data with custom ContentType.
	Render(http.ResponseWriter) error
	// WriteContentType writes custom ContentType.
	WriteContentType(w http.ResponseWriter)
}

// 验证类是否实现了接口
var (
	_ Render     = JSON{}
	_ Render     = IndentedJSON{}
	_ Render     = SecureJSON{}
	_ Render     = JsonpJSON{}
	_ Render     = AsciiJSON{}
	_ Render     = Data{}
	_ Render     = HTML{}
	_ HTMLRender = HTMLDebug{}
	_ HTMLRender = HTMLProduction{}
	_ Render     = Reader{}
	_ Render     = ProtoBuf{}
	_ Render     = Redirect{}
	_ Render     = String{}
	_ Render     = XML{}
	_ Render     = YAML{}
)

/**
 * 如果响应头部没有设置Content-Type，则设置响应头部中的Content-Type
 */
func writeContentType(w http.ResponseWriter, value []string) {
	// 获得响应的头部
	header := w.Header()
	// 设置响应头部中的Content-Type
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}
