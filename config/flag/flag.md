# NewFlagSet

创建独立参数集，生成一个与全局默认的 CommandLine（flag 包的全局 FlagSet）分离的参数集合，避免参数命名冲突或相互干扰

* flag.ContinueOnError：解析出错时返回错误，由调用方手动处理（如记录日志或重试）。
* flag.ExitOnError（常用）：出错时直接调用 os.Exit(2) 终止程序。
* flag.PanicOnError：触发 panic，适合需要堆栈追踪的调试场景。

```
flags = flag.NewFlagSet("goose", flag.ExitOnError)
```

# 参数类型
## String
```go
dir   = flags.String("dir", ".", "directory with migration files")
```

## Int
```go
startPort := startCmd.Int("port", 8080, "启动端口")
```

# usage
```go
flags.Usage = usage

func usage() {
	fmt.Println(usagePrefix)
	flags.PrintDefaults()
	fmt.Println(usageCommands)
}
````

# Parse
```go
if err := flags.Parse(os.Args[1:]); err != nil {
	log.Fatalf("goose: failed to parse flags: %v", err)
}
args := flags.Args()
if len(args) == 0 || args[0] == "-h" || args[0] == "--help" {
	flags.Usage()
	return
}
if len(args) < 3 {
	flags.Usage()
	return
}
```
