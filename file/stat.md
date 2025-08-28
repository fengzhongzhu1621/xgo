package file


# os.Stat()
目的： 此函数检索符号链接指向的文件或目录的信息。如果该文件是符号链接，os.Stat() 会跟踪到目标并检索目标文件的信息。
用法： 当您需要了解符号链接指向的实际文件详细信息时，请使用 os.Stat()。

# os.Lstat()
目的： 此函数检索有关符号链接本身的信息，而不跟踪链接。它返回有关符号链接本身的详细信息，如文件大小、权限和模式。
用法： 当您需要有关符号链接本身的信息时（例如，判断文件是否为符号链接），请使用 os.Lstat()。

# filepath.Walk()
当使用 filepath.Walk() 递归遍历目录时，务必谨慎处理符号链接，以避免无限循环或意外行为（例如，跟踪指向树中更高目录的符号链接）。
在这种情况下，使用 os.Lstat() 可以确保您不会在必要时跟踪符号链接。

# 错误码

* os.ErrNotExist 文件未找到： 如果路径不存在，这两个函数都将返回错误。
* os.ErrPermission 权限被拒绝： 如果程序没有权限访问文件或目录。
* 损坏的符号链接： 损坏的符号链接（指向不存在的文件的链接）将导致 os.Stat() 返回错误，但 os.Lstat() 将成功，返回有关符号链接本身的信息。

# Mode()&os.ModeSymlink
要检查文件是否为符号链接，请使用 os.Lstat() 并检查文件信息对象的 Mode()。可以验证 Mode()&os.ModeSymlink 是否为真，这表明该文件为符号链接。

```go
info, err := os.Lstat("example_symlink")
if err != nil {
    fmt.Println("Error:", err)
} else if info.Mode()&os.ModeSymlink != 0 {
    fmt.Println("This is a symbolic link")
}
```
