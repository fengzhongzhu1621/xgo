# 差异
与 flag 包不同的是，pflag 包参数定界符是两个 -，而不是一个 -，在 pflag 中 -- 和 - 具有不同含义
布尔类型的标志指定参数 --boolVar=false 需要使用等号 = 而非空格

# NewFlagSet
```go
fs1 := pflag.NewFlagSet("fs1", pflag.ExitOnError)
fs1.String("flag1", "default1", "This is flag1 from fs1")

flagset := pflag.NewFlagSet("test", pflag.ExitOnError)
flagset.SortFlags = false
flagset.Parse(os.Args[1:])

fmt.Printf("ip: %d\n", *ip)
fmt.Printf("boolVar: %t\n", boolVar)
fmt.Printf("host: %+v\n", h)

i, err := flagset.GetInt("ip")
fmt.Printf("i: %d, err: %v\n", i, err)
```
创建了一个名为 "fs1" 的 pflag.FlagSet，并定义了 flag1 参数。pflag 支持 POSIX/GNU 风格的命令行参数（如 --flag），且默认错误处理模式为 pflag.ExitOnError（解析失败时退出程序）

# AddGoFlagSet
将标准库的 fs2 合并到 pflag 的 fs1 中，使得 fs1 能解析 fs2 定义的参数（如 flag2）。这是 pflag 兼容标准库的关键方法
```go
fs1.AddGoFlagSet(fs2)
```

# Args
```go
fmt.Printf("NFlag: %v\n", pflag.NFlag()) // 返回已设置的命令行标志个数
fmt.Printf("NArg: %v\n", pflag.NArg())   // 返回处理完标志后剩余的参数个数
fmt.Printf("Args: %v\n", pflag.Args())   // 返回处理完标志后剩余的参数列表
fmt.Printf("Arg(1): %v\n", pflag.Arg(1)) // 返回处理完标志后剩余的参数列表中第 i 项
```

# SortFlags
```go
flagset.SortFlags = false
```

# Parse
```go
fs1.Parse([]string{"--flag1=value1", "--flag2=value2"})
pflag.Parse()
flagset.Parse(os.Args[1:])
```

# Lookup
```go
fmt.Println("flag1:", fs1.Lookup("flag1").Value.String())
```

# PrintDefaults
```go
pflag.PrintDefaults()
```

# 参数

## String
仅支持长格式（如 --name），需指定参数名称、默认值和帮助信息，返回一个指向字符串的指针

返回 *string 类型的指针，需通过解引用获取实际值

```go
fs1.String("flag1", "default1", "This is flag1 from fs1")
// 需指定参数名称、默认值和帮助信息，返回一个指向字符串的指针
// 示例调用：--name=Alice
var name = pflag.String("name", "default", "Description of the flag")
fmt.Println(*name) // 输出解析后的值
```

## StringP
支持长格式和短格式（如 --name 和 -n），需额外指定短选项的单个字符
```go
// 示例调用：-n Alice 或 --name=Alice
var name = pflag.StringP("name", "n", "default", "Description of the flag")
fmt.Println(*name) // 输出解析后的值
```

## Bool
```go
version := cmdline.Bool("version", false, "show version information")
if *ver {
	version.ShowVersion()
	os.Exit(0)
}
```

## BoolVarP
```go
var help bool
pflag.CommandLine.BoolVarP(&help, "help", "h", false, "show help info")

var boolVar bool
flagset.BoolVarP(&boolVar, "boolVar", "b", true, "help message for boolVar")
```

## Int
```go
var ip *int = pflag.Int("ip", 1234, "help message for ip")
fmt.Printf("ip: %d\n", *ip)
```

## IntP
```go
var ip = flagset.IntP("ip", "i", 1234, "help message for ip")
```

## IntVar
```go
var port int
pflag.IntVar(&port, "port", 8080, "help message for port")
fmt.Printf("port: %d\n", port)
```

## Var
```go
type host struct {
  value string
}

func (h *host) String() string {
  return h.value
}

func (h *host) Set(v string) error {
  h.value = v
  return nil
}

func (h *host) Type() string {
  return "host"
}

var h host
pflag.Var(&h, "host", "help message for host")
fmt.Printf("host: %+v\n", h)

var h host
flagset.VarP(&h, "host", "H", "help message for host")
```
