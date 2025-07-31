GORM 提供了一些默认的序列化器：json、gob、unixtime

```go
type User struct {
    Name        []byte                 `gorm:"serializer:json"`
    Roles       Roles                  `gorm:"serializer:json"`
    Contracts   map[string]interface{} `gorm:"serializer:json"`
    JobInfo     Job                    `gorm:"type:bytes;serializer:gob"`
    CreatedTime int64                  `gorm:"serializer:unixtime;type:time"` // 将 int 作为日期时间存储到数据库中
}

type Roles []string

type Job struct {
    Title    string
    Location string
    IsIntern bool
}
```

# 注册序列化器
```go
import "gorm.io/gorm/schema"

type SerializerInterface interface {
    Scan(ctx context.Context, field *schema.Field, dst reflect.Value, dbValue interface{}) error
    SerializerValuerInterface
}

type SerializerValuerInterface interface {
    Value(ctx context.Context, field *schema.Field, dst reflect.Value, fieldValue interface{}) (interface{}, error)
}
```
