
# 1. 安装 swag 命令
```go
go install github.com/swaggo/swag/cmd/swag
```

# 2. 初始化
```sh
swag init

# 根据指定的模板生成
swag init --template template/custom.tmpl

# 生成 YAML 格式的文档
swag init --output docs --format yaml
```

swag 会扫描你的代码并生成文档文件，默认会在 docs/ 文件夹下生成 docs.go 和 swagger.json

http://127.0.0.1:8000/swagger/index.html

注意：
1. swag init一定要和 main.go 处于同一级目录
2. main程序名称必须为 main.go

# 3. 通用API注释

## 3.1 基本注释
* @title：文档的标题。
* @version：API 的版本。
* @description：对 API 的描述。
* @termsOfService：服务条款的 URL。
* @contact.name、**@contact.url、@contact.email**：联系信息。
* @host：API 的主机地址。
* @BasePath：API 的基础路径。

### 3.1.1 定义多版本 API

```go
// @BasePath /api/v1
// @BasePath /api/v2
```


## 3.2 路由注释
* @Router：指定接口的 URL 路径和 HTTP 方法。
* @Summary：接口的简要描述。
* @Description：接口的详细描述。
* @Tags：接口的分类标签。
* @Accept：指定请求的 MIME 类型，例如 json。
* @Produce：指定响应的 MIME 类型，例如 json。
* @Param：描述接口的参数，例如查询参数、路径参数、请求体等。
* @Success：描述接口的成功响应。
* @Failure：描述接口的失败响应。
* @Header：描述响应头。
* @Security：描述接口的安全性，例如 BasicAuth、APIKey、BearerToken 等。

### 3.2.1 参数注释
@Param 用于描述接口的请求参数，支持路径参数、查询参数、请求体等多种类型。

参数类型
* path：路径参数。
* query：查询参数。
* body：请求体

值类型
* int
* string
* bool
* object
* array
* file


```go
// @Param id path int true "用户ID"
// @Param name query string false "用户名"
// @Param body body models.User true "用户信息"
```

```go
// @Param user body User true "用户信息" default({"name": "John Doe", "email": "john@example.com"})
// @Param user body User true "用户信息" example({"name": "John Doe", "email": "john@example.com"})
```

```go
// @Param status query string true "订单状态" Enums(pending, paid, shipped)
```

### 3.2.2 响应注释
@Success 和 @Failure 用于描述接口的成功和失败响应。

```go
// @Success 200 {object} models.User "成功时返回的用户信息"
// @Failure 400 {string} string "请求参数错误"
// @Failure 404 {string} string "找不到指定的资源"
```

```go
// @Success 200 {object} User "创建成功返回的用户信息" example({"id": 1, "name": "John Doe", "email": "john@example.com"})
```

### 3.2.3 安全性
@Security 用于描述 API 的安全机制。
```go
// @Security ApiKeyAuth
// @Security BasicAuth
```

### 3.2.4

```go
// @Header 200 {string} X-Token "服务器返回的 Token"
```

# 4. 运行时动态设置 Swagger 信息

```go
import "github.com/swaggo/swag"

swag.SwaggerInfo.Title = "My Dynamic API"
swag.SwaggerInfo.Host = "api.example.com"
swag.SwaggerInfo.Version = "2.0"
swag.SwaggerInfo.BasePath = "/v2"
```

# 5. 注解设置全局文档信息

通过在 main.go 文件中使用注解设置全局文档信息。

```go
// @title My API
// @version 1.0
// @description 这是我的 API 文档。
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /api/v1
```

# 6. 自定义模板

创建一个自定义模板文件，例如 template/custom.tmpl

```sh
swag init --template template/custom.tmpl
```

```json
{
  "swagger": "2.0",
  "info": {
    "title": "{{.Title}}",
    "description": "{{.Description}}",
    "version": "{{.Version}}"
  },
  "paths": {
    {{range .Paths}}
    "{{.Path}}": {
      "get": {
        "summary": "{{.Summary}}",
        "operationId": "{{.OperationID}}",
        "parameters": [{{range .Parameters}}{{.Name}}: {{.Type}}{{end}}]
      }
    }
    {{end}}
  }
}
```

# 7. 定义全局参数

通过注释为文档添加全局参数

通过 @securityDefinitions 注解定义一个全局的 API 密钥认证方式，所有需要此认证的路由都可以共享这个定义。

```go
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
````

# 8. 生成离线文档
swaggo/swag 生成的 Swagger 文档（swagger.json 或 swagger.yaml）可以通过 swagger-ui 工具生成离线文档，方便你在无网络环境下使用。
使用 Swagger-UI 官方提供的工具，将生成的 swagger.json 或 swagger.yaml 嵌入到静态 HTML 文件中，形成离线可访问的文档。

```sh
swagger-ui-dist/swagger-ui-bundle.js swagger.json > index.html
```

# 示例
```go
// @Summary 获取用户信息
// @Description 根据用户ID获取用户详细信息
// @Tags user
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} User "返回用户信息"
// @Failure 404 {object} string "用户未找到"
// @Router /user/{id} [get]
func getUser(c *gin.Context) {
 // 实现获取用户的逻辑
}
```

通过在注释中使用 {object} User，swaggo/swag 会自动生成该结构体及其嵌套结构体的描述。
