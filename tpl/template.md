Templ 是 Go 语言自带的模板引擎，它提供了一种简单而强大的方式来生成动态HTML内容。

# 变量和表达式

使用 {{ 和 }} 包裹变量和表达式

```go
<h1>{{.Title}}</h1>
<p>{{.Content}}</p>
```

# ParseFiles

```go
tmpl, _ := template.ParseFiles("template.html")
data := struct {
    Title   string
    Content string
}{"Hello", "World"}
tmpl.Execute(w, data)
```

# 条件
```go
{{if .LoggedIn}}
<p>Welcome, {{.Username}}</p>
{{else}}
<a href="/login">Login</a>
{{end}}
```

# 循环

```go
{{range .Users}}
<li>{{.Name}}</li>
{{end}}
```

# 自定义函数和管道

```go
{{.Content | truncate 100}}
```