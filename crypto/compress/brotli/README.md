# 简介
Brotli 是由 Google 于 2015 年开源的一种通用无损压缩算法，它基于 LZ77 算法进行改进，并结合了上下文建模和熵编码等技术，能够在压缩率和速度之间取得更佳的平衡。Brotli 尤其擅长处理文本和网页内容，相比传统的 Gzip 算法，它能够实现更高的压缩率，从而有效减少文件大小，加快页面加载速度。

# 安装
```
brew install brotli
```

# 优势
* 更高的压缩率: 相比 Gzip，Brotli 平均能够将文本文件压缩 15% - 25%，对于网页内容，压缩率甚至可以达到 30% 以上。
* 更快的解压速度: 尽管 Brotli 的压缩速度略慢于 Gzip，但它的解压速度却与之相当，甚至更快。
* 广泛的浏览器支持: 目前，主流的浏览器都已支持 Brotli 算法，包括 Chrome、Firefox、Edge、Safari 等，这意味着使用 Brotli 压缩的网页内容能够被绝大多数用户直接访问。

# 对比
| 文件大小 | 压缩算法 | 压缩后大小 | 压缩率 | 压缩时间 | 解压时间 |
|----------|----------|------------|--------|----------|----------|
| 1 KB     | Gzip     | 0.5 KB     | 50%    | 0.1 ms   | 0.05 ms  |
| 1 KB     | Brotli   | 0.4 KB     | 60%    | 0.15 ms  | 0.06 ms  |
| 10 KB    | Gzip     | 4 KB       | 60%    | 0.2 ms   | 0.1 ms   |
| 10 KB    | Brotli   | 3 KB       | 70%    | 0.3 ms   | 0.12 ms  |
| 100 KB   | Gzip     | 30 KB      | 70%    | 1 ms     | 0.5 ms   |
| 100 KB   | Brotli   | 20 KB      | 80%    | 1.5 ms   | 0.6 ms   |
| 1 MB     | Gzip     | 250 KB     | 75%    | 8 ms     | 4 ms     |
| 1 MB     | Brotli   | 180 KB     | 82%    | 12 ms    | 4.5 ms   |

Brotli 在压缩率方面 consistently 优于 Gzip，尤其是在处理包含重复模式的大型文件时，优势更加明显。尽管 Brotli 的压缩时间略长，但其解压速度与 Gzip 相当，而且更高的压缩率使其成为更优的选择。

## 压缩算法对比结论

### 1. 压缩率对比
| 文件大小 | Gzip压缩率 | Brotli压缩率 | Brotli优势 |
|----------|------------|--------------|------------|
| 1 KB     | 50%        | 60%          | +10%       |
| 10 KB    | 60%        | 70%          | +10%       |
| 100 KB   | 70%        | 80%          | +10%       |
| 1 MB     | 75%        | 82%          | +7%        |

**结论**：
Brotli在所有测试场景下都表现出更高的压缩率，平均比Gzip高约8-10%，特别是在小文件（1KB-100KB）上优势更明显。

---

### 2. 压缩速度对比
| 文件大小 | Gzip压缩时间 | Brotli压缩时间 | Gzip优势 |
|----------|--------------|----------------|----------|
| 1 KB     | 0.1 ms       | 0.15 ms        | -0.05 ms |
| 10 KB    | 0.2 ms       | 0.3 ms         | -0.1 ms  |
| 100 KB   | 1 ms         | 1.5 ms         | -0.5 ms  |
| 1 MB     | 8 ms         | 12 ms          | -4 ms    |

**结论**：
Gzip的压缩速度始终快于Brotli，特别是在大文件（1MB）上优势显著（快33%）。Brotli的压缩时间随文件大小增长更快。

---

### 3. 解压速度对比
| 文件大小 | Gzip解压时间 | Brotli解压时间 | Gzip优势 |
|----------|--------------|----------------|----------|
| 1 KB     | 0.05 ms      | 0.06 ms        | -0.01 ms |
| 10 KB    | 0.1 ms       | 0.12 ms        | -0.02 ms |
| 100 KB   | 0.5 ms       | 0.6 ms         | -0.1 ms  |
| 1 MB     | 4 ms         | 4.5 ms         | -0.5 ms  |

**结论**：
Gzip的解压速度也普遍快于Brotli，但差距比压缩速度小（约10-20%），且在小文件上差异不明显。

---

### 4. 综合建议
1. **优先选择Brotli的场景**：
   - 对存储空间敏感（如Web资源、移动端应用）
   - 小文件压缩（<100KB）
   - 需要更高压缩率的场景

2. **优先选择Gzip的场景**：
   - 对压缩速度敏感（如实时数据传输）
   - 大文件压缩（>100KB）
   - 兼容性要求高的环境（Brotli支持较新）

3. **折中方案**：
   - 可根据文件大小动态选择算法（小文件用Brotli，大文件用Gzip）
   - Web服务可同时提供两种格式供客户端选择（通过`Accept-Encoding`头）

---

### 5. 关键发现
- **压缩率**：Brotli完胜（+7-10%）
- **速度**：Gzip完胜（特别是大文件）
- **权衡**：Brotli牺牲部分速度换取更高压缩率，适合对存储敏感的场景
