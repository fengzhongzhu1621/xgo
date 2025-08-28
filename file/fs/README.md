# FS

```go
type FS interface {
    Open(name string) (File, error)
}
```

# File

```go
type File interface {
    Stat() (FileInfo, error)
    Read([]byte) (int, error)
    Close() error
}
```
