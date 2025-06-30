# 引导套件
go install github.com/onsi/ginkgo/v2/ginkgo@latest

```
. "github.com/onsi/ginkgo/v2"
. "github.com/onsi/gomega"
```

```
ginkgo bootstrap
```

生成一个名为 books_suite_test.go 的文件

# 添加 Specs 到你的套件

```
ginkgo generate book
```

# Describe
Describe块用于组织Specs，将一个或多个测试例归类

Ordered: 测试套件中的测试用例将按照它们在代码中出现的顺序执行
```go
var _ = Describe("Account", Ordered, func() {
```

## BeforeEach
会在每个测试用例（It 语句）执行之前运行。它通常用于设置测试环境，例如初始化变量、创建对象等。如果有多个 BeforeEach 函数，它们会按照声明的顺序依次执行

## AfterEach
每个测试例执行后执行该段代码

# It / Specify
是测试例的基本单位，即It包含的代码就算一个测试用例
Specify和It功能完全一样，It属于其简写

```go
	Describe("delete app api", func() {
		It("should delete app permanently", func() {...})
		It("should delete app failed if services existed", func() {...})
```

## JustBeforeEach
在BeforeEach执行之后，测试例执行之前执行。将创建与配置分离。通常用于执行一些在测试用例之间不需要重复的操作，例如记录测试结果、更新统计信息等。如果有多个 JustBeforeEach 函数，它们会按照声明的顺序依次执行。


# BeforeSuite
是在该测试集执行前执行，即该文件夹内的测试例执行之前

# AfterSuite
是在该测试集执行后执行，即该文件夹内的测试例执行完后

# Context
将一个或多个测试例归类，丰富Describe所描述的行为或方法，增加条件语句，尽可能全地覆盖各种condition

* Describe 用于描述你的代码的一个行为
* Context 用于区分上述行为的不同情况，通常为参数不同导致

# Measure
用于性能测试

# By
是打印信息，内容只能是字符串，只会在测试例失败后打印，一般用于调试和定位问题

# Fail
是标志该测试例运行结果为失败，并打印里面的信息


# 执行测试

```
# 运行当前目录中的测试
ginkgo

# 运行其它目录中的测试
ginkgo /path/to/package /path/to/other/package ...

# 递归运行所有子目录中的测试
ginkgo -r ...

# 传递参数给测试套件
ginkgo -- PASS-THROUGHS-ARGS

# 跳过某些包
ginkgo -skipPackage=PACKAGES,TO,SKIP

# 启用并行测试，Ginkgo会自动创建适当数量的节点（进程）。你也可以指定节点数量
ginkgo -p
ginkgo -nodes=N

ginkgo -v

ginkgo --json-report ./ginkgo.report -focus "String"   -r
```


# gomega

## 构建断言
### To
.To() 方法是 gomega 断言的原始形式，它遵循了行为驱动开发（BDD）的一些原则，特别是在断言语法上。使用 .To() 方法时，断言读起来更像是自然语言的描述，这有助于非技术人员理解测试的目的。
.To() 后面通常会跟一个匹配器（Matcher），用于定义期望的条件。.To() 方法本身并不执行任何比较操作，而是将实际的值和期望的条件传递给匹配器进行检查。

.To() 和 .Should() 在功能上是等价的

```
Expect(book.Title).To(Equal("Les Miserables"))
```

### NotTo
```
Expect(ACTUAL).NotTo(Equal(EXPECTED))
```

### ToNot
```
Expect(ACTUAL).ToNot(Equal(EXPECTED))
```

### Should
为了提供另一种更简洁的断言语法而引入的。它去掉了 .To() 前缀，使得断言语句更加紧凑。
```
Expect(book.Title).Should(Equal("Les Miserables"))
Expect(err).Should(gomega.BeNil())
Expect(Validate(inputData[1])).Should(gomega.MatchError("只要男的"))
```

## 匹配器
### Equal
```
Expect(book.Title).To(Equal("Les Miserables"))
```

### BeNil
```
Expect(err).Should(gomega.BeNil())
```
### MatchError
```
Expect(Validate(inputData[1])).Should(gomega.MatchError("只要男的"))
```
