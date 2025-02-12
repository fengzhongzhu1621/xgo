# 简介
bcrypt是一个由美国计算机科学家尼尔斯·普罗沃斯（Niels Provos）以及大卫·马齐耶（David Mazières）根据Blowfish加密算法所设计的密码散列函数，
于1999年在USENIX中展示。实现中bcrypt会使用一个加盐的流程以防御彩虹表攻击，同时bcrypt还是适应性函数，
它可以借由增加迭代之次数来抵御日益增进的电脑运算能力透过暴力法破解。

由bcrypt加密的文件可在所有支持的操作系统和处理器上进行转移。它的口令必须是8至56个字符，并将在内部被转化为448位的密钥。
然而，所提供的所有字符都具有十分重要的意义。密码越强大，数据就越安全。

除了对数据进行加密，默认情况下，bcrypt在删除数据之前将使用随机数据三次覆盖原始输入文件，以阻挠可能会获得计算机数据的人恢复数据的尝试。
如果您不想使用此功能，可设置禁用此功能。

# 不可逆性
bcrypt就是一种加盐的单向Hash，不可逆的加密算法，同一种明文（plaintext），每次加密后的密文都不一样，而且不可反向破解生成明文，破解难度很大。
生成的密文是60位的，比MD5更安全，但加密更慢。


# cost
```golang
const (
        // 传递给GenerateFromPassword的最小允许开销
	MinCost     int = 4
        // 传递给GenerateFromPassword的最大允许开销
	MaxCost     int = 31
        // 如果将低于MinCost的cost传递给GenerateFromPassword，则实际设置的cost
	DefaultCost int = 10
)
```
