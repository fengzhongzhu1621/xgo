package tpl

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/http"
	"text/template"
)

//go:embed index.html
var indexTmpl string

type Person struct {
	ID   int // ID 对应 html 模板中的{{.ID}}占位符变量，严格区分大小写，其他几个字段类似
	Name string
	Age  int
	City string
}

func main() {
	// 注册路由 index，用于渲染 web 页面
	http.HandleFunc("/index", func(writer http.ResponseWriter, request *http.Request) {
		// 模拟数据
		people := []Person{
			{1, "John Doe", 30, "New York"},
			{2, "Jane Smith", 25, "Los Angeles"},
			{3, "Mary Johnson", 35, "Chicago"},
		}

		// 读取 index.html内容并解析模板
		t := template.Must(template.New("webpage").Parse(indexTmpl))
		// 创建一个缓冲区，用于存储渲染后的结果
		renderWriter := bytes.NewBuffer(nil)
		// 执行模板渲染，将渲染后的结果输出到 Buffer
		err := t.Execute(renderWriter, people)
		if err != nil { // 渲染失败，返回错误信息
			writer.WriteHeader(http.StatusBadRequest)
			_, _ = writer.Write([]byte(err.Error()))
			return
		}

		// 渲染成功，返回渲染后的结果
		writer.WriteHeader(http.StatusOK)
		writer.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = writer.Write(renderWriter.Bytes())
	})

	// 启动服务
	// http://localhost:8080/index
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println("Starting server error:", err)
	}
}
