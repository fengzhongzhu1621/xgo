# 简介
用于 YAML 和 shell 命令模板，指的是对标 text/template 的安全增强。是 2023 年初发布的第一个安全库的家族成员。
内部在使用 text/template 做开发基于 YAML 的应用程序时，经常受到 YAML 注入的攻击。当检测到注入时，SafeText 库会返回错误。

* 防止命令注入
* 安全的变量替换
* 严格的语法检查
* 内置转义机制

## yamltemplate
这个库专注于安全地生成YAML配置文件。它默认阻止输入字符串影响YAML结构，仅允许它们作为值存在。
对于需要改变结构的部分，需要明确使用StructuralData注解，并可能需要进一步的数据验证。

## shtemplate
设计用于生成安全的shell脚本。默认情况下，它不允许输入字符串插入新的命令或标志，而需要通过StructuralData注解来显式启用。
此外，它还有一个AllowFlags注解，用于允许输入作为命令标志。

## shsprintf
这是一个类似于fmt.Sprintf的安全版本，专为shell脚本设计。
它保证了即使在不正确的转义情况下，也不会发生命令或标志注入。
