# 简介

防止路径遍历攻击（又称 zip slip）以及与处理归档文件相关的各种攻击。
可直接替换 Go 标准库中的的 archive/tar 和 archive/zip，直接换包的导入路径就可以了。使用后，压缩包中如果包含恶意信息，发现后将会被清除。

* 防止 ZIP/TAR 炸弹
* 路径穿越保护
* 文件大小限制
* 符号链接控制
* 文件类型验证
