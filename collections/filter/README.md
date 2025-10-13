## 原理

理解拦截器的原理关键点在于理解拦截器的`触发时机` 以及 `顺序性`。

触发时机：拦截器可以拦截到接口的请求和响应，并对请求，响应，上下文进行处理（用通俗的语言阐述也就是 可以在`请求接受前`做一些事情，`请求处理后`做一些事情），因此，拦截器从功能上说是分为两个部分的 前置（业务逻辑处理前） 和 后置（业务逻辑处理后）

## FAQ

### Q：拦截器入口这里能否拿到二进制数据

不可以，拦截器入口这里的 req rsp 都是已经经过序列化过的结构体了，可以直接使用数据，没有二进制。


### Q：多个拦截器执行顺序如何

多个拦截器的执行顺序按配置文件中的数组顺序执行，如

```yaml
server:
  filter:
    - filter1
    - filter2
  service:
    - name: trpc.app.server.service
      filter:
        - filter3
```

则执行顺序如下：

```
接收到请求 -> filter1 前置逻辑 -> filter2 前置逻辑 -> filter3 前置逻辑 -> 用户的业务处理逻辑 -> filter3 后置逻辑 -> filter2 后置逻辑 -> filter1 后置逻辑 -> 回包
```

### Q：一个拦截器必须同时设置 server 和 client 吗

不需要，只需要 server 时，client 传入 nil，同理只需要 client 时，server 传入 nil，如

```golang
filter.Register("name1", serverFilter, nil)  // 注意，此时的 name1 拦截器只能配置在 server 的 filter 列表里面，配置到 client 里面，rpc 请求会报错
filter.Register("name2", nil, clientFilter)  // 注意，此时的 name2 拦截器只能配置在 client 的 filter 列表里面，配置到 server 里面会启动失败
```
