
并不是所有的 异常都能被 recover 捕获

https://stackoverflow.com/questions/57486620/are-all-runtime-errors-recoverable-in-go


如果要 recover 异常，那么至少这个异常是原子可恢复的，不能在 recover 之后留一个不可用的内存环境给到后续的业务逻辑代码。下面是几类不能被捕获的异常。

* Out of memory
* Concurrent map writes and reads
* Stack memory exhaustion
* Attempting to launch a nil function as a goroutine
* All goroutines are asleep - deadlock
* Thread limit exhaustion
