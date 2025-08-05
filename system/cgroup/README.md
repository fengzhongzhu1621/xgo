# 内存
* Total: cgroup 被限制可以使用多少内存，可以从文件里的 hierarchical_memory_limit 获得，但不是所有 cgroup 都限制内存，没有限制的话会获得 2^64-1 这样的值，我们还需要从 /proc/meminfo 中获得 MemTotal，取两者最小。
* RSS: Resident Set Size 实际物理内存使用量，在 memory/memory.stat 的 rss 只是 anonymous and swap cache memory，文档里也说了如果要获得真正的 RSS 还需要加上 mapped_file。
* Cached: memory/memory.stat 中的 cache
* MappedFile: memory/memory.stat 中的 mapped_file

这里不将共享内存计算入容器使用内存
