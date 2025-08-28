package httpstaticserver

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/fengzhongzhu1621/xgo/crypto/compress/zipfile"
	"github.com/fengzhongzhu1621/xgo/validator"
)

func (s *HTTPStaticServer) hUploadOrMkdir(w http.ResponseWriter, req *http.Request) {
	// 获得需要上传的和服务器目录地址
	dirpath := s.getRealPath(req)

	// 判断用户是否有上传文件的权限
	auth := s.readAccessConf(dirpath)
	if !auth.CanUpload(req) {
		http.Error(w, "Upload forbidden", http.StatusForbidden)
		return
	}

	// 获得用户上传的文件
	f, header, err := req.FormFile("file")

	// 目录不存在则创建
	if _, err := os.Stat(dirpath); os.IsNotExist(err) {
		if err := os.MkdirAll(dirpath, os.ModePerm); err != nil {
			log.Println("Create directory:", err)
			http.Error(w, "Directory create "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 如果上传的附件不存在，则仅仅创建目录
	if f == nil { // only mkdir
		w.Header().Set("Content-Type", "application/json;charset=utf-8")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     true,
			"destination": dirpath,
		})
		return
	}

	if err != nil {
		log.Println("Parse form file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer func() {
		f.Close()
		req.MultipartForm.RemoveAll() // Seen from go source code, req.MultipartForm not nil after call FormFile(..)
	}()

	// 获得上传的文件名
	filename := req.FormValue("filename")
	if filename == "" {
		filename = header.Filename
	}
	// 判读文件中是否存在特殊字符
	if err := validator.CheckFilename(filename); err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	// 创建空文件
	dstPath := filepath.Join(dirpath, filename)
	var copyErr error
	dst, err := os.Create(dstPath)
	if err != nil {
		log.Println("Create file:", err)
		http.Error(w, "File create "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Note: very large size file might cause poor performance
	// _, copyErr = io.Copy(dst, file)
	// 从 sync.Pool 中获取一个对象。如果池中有可用的对象，则返回该对象；否则，调用用户提供的 New 函数来创建一个新对象。
	// 第一次调用 Get 时，由于池中没有可用的对象，因此会调用 NewMyStruct 函数来创建一个新对象。
	// 然后，我们使用 pool.Put(obj) 将对象放回池中。
	// 再次调用 pool.Get() 方法从池中获取对象。这次，由于池中有可用的对象（之前放回的对象），因此会复用这个对象，而不是创建一个新对象。
	buf := s.bufPool.Get().([]byte)
	//  通过 Put 复用 buf对象
	defer s.bufPool.Put(buf)
	// 用于从一个 io.Reader 接口复制数据到另一个 io.Writer 接口。与 io.Copy 不同，io.CopyBuffer 允许你指定一个自定义的缓冲区来执行复制操作。
	// 这在处理大量数据时可能更高效，因为你可以控制缓冲区的大小。
	_, copyErr = io.CopyBuffer(dst, f, buf)
	dst.Close()
	// }
	if copyErr != nil {
		log.Println("Handle upload file:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")

	// 解压缩文件
	if req.FormValue("unzip") == "true" {
		err = zipfile.UnzipFile(dstPath, dirpath)
		os.Remove(dstPath)
		message := "success"
		if err != nil {
			message = err.Error()
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success":     err == nil,
			"description": message,
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":     true,
		"destination": dstPath,
	})
}
