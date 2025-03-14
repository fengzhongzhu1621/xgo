# cobra

cobra是一个非常流行的库，用于创建强大的现代CLI（命令行界面）应用程序。

```
go get -u github.com/spf13/cobra/cobra
```

## StringVarP

```go
serveCmd.Flags().StringVarP(&cfgFile, "config", "c", "", "config file (default is config.yml;required)")
```