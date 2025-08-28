# 简介
如果项目中大量依赖了Redis 存储，可以使用上文中gomock 对Redis conn接口进行mock。

也可以在本地起一个miniredis 服务，它实现了大部分Redis server功能，便于本地单元测试的执行

```sh
go get github.com/alicebob/miniredis/v2
```
