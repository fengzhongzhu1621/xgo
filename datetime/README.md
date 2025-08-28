# timedatectl
```bash
$ timedatectl
      Local time: Mon 2024-10-28 09:41:26 CST
  Universal time: Mon 2024-10-28 01:41:26 UTC
        RTC time: Mon 2024-10-28 01:41:26
       Time zone: Asia/Shanghai (CST, +0800)
     NTP enabled: no
NTP synchronized: no
 RTC in local TZ: no
      DST active: n/a
```

列出所有时区
```bash
timedatectl list-timezones
```

# 硬件时钟（Hardware Clock）
硬件时钟通常被称为硬件时钟、实时时钟（RTC）、BIOS 时钟和 CMOS 时钟。Linux 内核也称其为持久时钟（persistent clock）。
具有自己的电源域（如电池、电容等），即使在机器关闭或断电时也能运行。硬件时钟的基本目的是在 Linux 不运行时保持时间，以便在启动时从硬件时钟初始化系统时钟。

# 系统时钟（System Clock）
系统时钟是 Linux 内核的一部分，由定时器中断驱动。在 ISA 机器上，定时器中断是 ISA 标准的一部分。系统时间是自 1970 年 1 月 1 日 00:00:00 UTC 来的秒数。
系统时间不是整数，实际上具有无限精度。

## 内核时区
Linux 内核的时区由 hwclock 设置。

* tz_minuteswest：表示本地时间（未调整夏令时）比 UTC 落后多少分钟。
* tz_dsttime：表示当前本地使用的夏令时（DST）类型。这个字段在 Linux 下不使用，总是为零。

