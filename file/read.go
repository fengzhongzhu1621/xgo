package file

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"strconv"
	"strings"
)

// StatFile 判断文件类型并返回一个 io.ReadCloser，用于后续读取文件内容
func StatFile(filename string) (info os.FileInfo, reader io.ReadCloser, err error) {
	// 用于获取文件或目录的元信息（metadata），但不会跟随符号链接
	info, err = os.Lstat(filename)
	if err != nil {
		return info, nil, err
	}
	// 用于判断文件或目录是否为符号链接
	if info.Mode()&os.ModeSymlink != 0 {
		var target string
		// 读取符号链接指向的目标路径
		// 请注意，在使用 os.Readlink 之前，确保提供的路径确实是一个符号链接。否则，可能会遇到错误。
		// 可以使用 os.Lstat 函数和 FileInfo.Mode().IsSymlink() 方法来检查路径是否为符号链接。
		target, err = os.Readlink(filename)
		if err != nil {
			return info, nil, err
		}
		// 实现了 io.Closer 接口，但不会执行任何实际操作
		// 可以将任何实现了 io.Reader 接口的类型转换为 io.Closer 接口类型，而无需实际实现关闭操作。
		// bytes.NewBuffer 接受一个字节切片（[]byte）作为参数，并返回一个包含这些字节的缓冲区（*bytes.Buffer 类型）
		reader = io.NopCloser(bytes.NewBuffer([]byte(target)))
	} else if !info.IsDir() {
		// 打开文件
		reader, err = os.Open(filename)
		if err != nil {
			return info, reader, err
		}
	} else {
		reader = io.NopCloser(bytes.NewBuffer(nil))
	}

	return info, reader, err
}

// ReadFileContent 读取文件的内容
func ReadFileContent(filepath string) (string, error) {
	// 用于读取指定文件的全部内容，并将其作为一个字节切片返回。如果读取过程中发生错误，函数会返回一个非空的错误值。
	// os.ReadFile 函数会一次性读取整个文件内容到内存中，因此对于非常大的文件，可能会导致内存不足的问题。
	// 在这种情况下，建议使用 os.Open 和 io.Reader 接口逐块读取文件内容。
	content, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(content), nil
}

// ReadFromFile 从文件中读取内容
func ReadFromFile(name string) (ret string, err error) {
	file, err := os.Open(name)
	if err != nil {
		return ret, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return ret, err
	}

	return strings.TrimSpace(string(data)), nil
}

// ReadUint64FromFile 从文件中读取 uint64 类型内容
func ReadUint64FromFile(file string) (ret uint64, err error) {
	dat, err := ReadFromFile(file)
	if err != nil {
		return ret, err
	}

	ret, err = strconv.ParseUint(dat, 10, 64)
	return
}

// ReadInt64FromFile 从文件中度一个 int64 类型的数
func ReadInt64FromFile(file string) (ret int64, err error) {
	dat, err := ReadFromFile(file)
	if err != nil {
		return ret, err
	}

	ret, err = strconv.ParseInt(dat, 10, 64)
	return
}

// ReadMapFromFile 从文件中读取一系列数据
// 从一个文本文件中读取数据，并将这些数据解析成一个 map[string]uint64（即字符串到无符号 64 位整数的映射）
func ReadMapFromFile(name string) (ret map[string]uint64, err error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// 存储最终解析出来的数据
	ret = make(map[string]uint64)

	// 创建一个扫描器，用于按行读取文件内容
	scanner := bufio.NewScanner(file)
	// 逐行读数据，每次调用会读取下一行，直到文件结束或出错。它返回一个布尔值，表示是否还有下一行。
	for scanner.Scan() {
		// 获取当前扫描到的行的文本内容（字符串）
		line := scanner.Text()
		// 假设每行的格式应该是：key value
		// 将当前行按第一个空格 " " 分割成最多 2 部分，结果存入 items 切片
		// 比如，若 line = "key 123"，则 items = ["key", "123"]
		items := strings.SplitN(line, " ", 2)
		if len(items) != 2 {
			continue
		}

		key := strings.TrimSpace(items[0])
		valueStr := strings.TrimSpace(items[1])

		if key == "" || valueStr == "" {
			continue // 跳过空 key 或空 value 的行
		}

		// 取分割后的第二部分，即我们期望是数字的部分，尝试将字符串解析为 uint64 类型的整数
		v, err := strconv.ParseUint(valueStr, 10, 64)
		if err != nil {
			continue
		}
		ret[key] = v
	}

	// 检测文件损坏或读取中断
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return ret, nil
}
