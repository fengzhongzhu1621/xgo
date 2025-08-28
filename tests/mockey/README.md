# mockey
Mockey是字节跳动开源的Golang单元测试Mock框架，基于gomock和goconvey封装，提供了更直观的接口和更便捷的使用方式，常用于单元测试中模拟函数和变量的行为，以隔离外部依赖，提升测试稳定性和效率。

```go
    . "github.com/bytedance/mockey"
    . "github.com/smartystreets/goconvey/convey"
```

# 对比
### Mockey vs. gomock vs. gomonkey vs. sqlmock/miniredis

| 对比项               | Mockey (字节跳动)                                                                 | gomock (Google)                                                                 | gomonkey (社区)                                                                 | sqlmock/miniredis (数据库专用)                                                |
|----------------------|----------------------------------------------------------------------------------|---------------------------------------------------------------------------------|---------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| **主要用途**         | 通用函数/方法/变量 Mock                                                           | 基于 Interface 的 Mock                                                          | 运行时函数/方法替换                                                             | 数据库操作 Mock                                                               |
| **Mock 对象类型**    | 函数、变量、方法（含私有方法）                                                    | 仅限 Interface 方法                                                             | 函数、方法                                                                      | 数据库驱动（如 `*sql.DB`）                                                    |
| **使用方式**         | 直接调用 `Mock()` + `When()` + `Return()`                                       | 需生成 Mock 代码（`mockgen`）                                                   | 直接调用 `ApplyFunc()` 等                                                       | 初始化模拟数据库（如 `miniredis.Run()`）                                      |
| **条件设置**         | 支持 `When()` 灵活条件                                                          | 通过 `EXPECT().Return()` 设置                                                  | 支持简单条件                                                                    | 不适用                                                                        |
| **返回值设置**       | 支持多返回值、序列（`sequence`）                                                 | 固定返回值                                                                      | 支持固定返回值                                                                  | 固定返回值                                                                    |
| **私有方法支持**     | ✅ 支持                                                                          | ❌ 不支持                                                                       | ✅ 支持                                                                         | ❌ 不支持                                                                     |
| **平台兼容性**       | macOS 无需高权限                                                                | macOS 可能需要高权限                                                            | macOS 无需高权限                                                                | 无特殊要求                                                                    |
| **测试组织**         | 提供 `PatchConvey`（类似 `goconvey`）                                           | 需手动组织 `Convey`                                                             | 需手动组织测试                                                                  | 需手动组织测试                                                                |
| **依赖注入**         | 不依赖 Interface                                                                | 强依赖 Interface                                                                | 不依赖 Interface                                                                | 不依赖 Interface                                                              |
| **适用场景**         | 复杂逻辑 Mock（如全局变量、私有方法）                                            | 标准 Interface 测试                                                             | 快速替换函数/方法                                                               | 数据库单元测试                                                                |

### 关键差异总结
1. **Mockey 的优势**：
   - **无需 Interface**：可直接 Mock 函数、变量、私有方法，适合遗留代码或非 Interface 场景。
   - **灵活的条件和返回值**：支持动态条件（`When`）和序列返回值（`sequence`）。
   - **平台友好**：macOS 无需高权限，适合本地开发。

2. **gomock 的优势**：
   - **类型安全**：基于 Interface 生成代码，编译时检查。
   - **Google 维护**：适合标准 Go 项目。

3. **gomonkey 的优势**：
   - **动态替换**：运行时修改函数/方法，适合无法通过 Interface Mock 的场景。

4. **sqlmock/miniredis 的优势**：
   - **数据库专用**：无需真实数据库，适合集成测试。

### 如何选择？
- **需要 Mock 函数/变量/私有方法** → **Mockey** 或 **gomonkey**
- **标准 Interface 测试** → **gomock**
- **数据库测试** → **sqlmock/miniredis**
