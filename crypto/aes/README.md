# 1. aes-everywhere

Cross Language AES 256 Encryption Library (Bash, Powershell, C#, Dart, GoLang, Java, JavaScript, Lua, PHP, Python, Ruby, Swift)

https://github.com/mervick/aes-everywhere
https://github.com/mervick/aes-everywhere/blob/master/go/aes256/aes256.go


# 2. AES-GCM
AES-GCM（Galois/Counter Mode）是一种使用AES算法的加密模式，它结合了Galois域乘法运算和计数器模式，提供加密、认证和完整性校验（MAC）功能。

## 特点

* **认证加密**：AES-GCM不仅提供数据的加密，还能验证数据的完整性和来源的真实性。
* **安全性**：通过使用额外的认证数据和消息认证码（MAC），AES-GCM能够抵御重放攻击和篡改攻击。
* **性能**：相比于其他加密模式，如CBC，GCM支持并行加密和解密，提高了加密和解密的效率

## 工作原理
* **加密过程**：AES-GCM使用一个初始向量（IV）和一个密钥来加密数据。在加密过程中，数据被分成固定大小的块，每个块都通过AES算法加密，并使用Galois域乘法运算生成一个MAC标签。
* **解密过程**：解密时，使用相同的密钥和IV对密文进行解密，并重新计算MAC标签。如果计算出的MAC标签与密文中包含的MAC标签匹配，则解密成功；否则，解密失败，表明数据在传输过程中可能被篡改

## 比较
* 与CBC模式的比较：CBC模式不提供消息的完整性校验，而GCM通过GMAC提供了这一功能。此外，GCM支持并行加密和解密，而CBC模式是串行执行的，因此在处理大量数据时，GCM的效率更高。
* 与CCM模式的比较：CCM模式也提供加密和认证功能，但它使用的是CBC-MAC而不是GCM中的GMAC。此外，CCM模式在某些情况下可能不如GCM安全。
