# os.Root
https://pkg.go.dev/os@master#Root

os.Root 可以锁定工作目录。 使用户无法打开目录外的文件，例如 ../../../etc/passwd 。 可以强制约束用户， 限制用户行为， 检查计划外的使用逻辑。

1. 使用 root, _ := os.OpenRoot(basedir) 锁定工作目录
2. 以后的所有操作都要基于 root.Xzzzz() 展开
3. root.OpenFile(path) 在打开文件之前， 会判断 文件路径 的合法性。 basedir 之外的路径拒绝访问。

   1. path 可以是绝对路径。
   2. path 也可以是基于 basedir 的相对路径。

不支持 root.ReadFile(path) 这样的快捷 API。 还是需要分成 root.OpenFile 和 io.ReadAll。
