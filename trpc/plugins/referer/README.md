# 简介
http referer 安全验证

# 使用
```go
import (
     _ "trpc.group/trpc-go/trpc-filter/referer"
)
```

# 配置
```yaml
server:
 ...
 filter:
  ...
  - referer
plugins:
  auth:
    referer:
      apply_to_all_path:
        - qq.com
      path1:
        - qq.com
      path2:
        - NULL
```

- apply_to_all_path 默认规则
- NULL 表示 referer 为空时可以放过
- `*` 匹配所有的(不安全)
- referer 插件 会对所有 url 验证 非白名单模式
