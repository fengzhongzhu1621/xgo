package file

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
)

func userHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	}
	return os.Getenv("HOME")
}

// 将路径转换为绝对路径.
func AbsPathify(inPath string) string {
	// 如果是home路径，则转换为绝对路径
	if inPath == "$HOME" || strings.HasPrefix(inPath, "$HOME"+string(os.PathSeparator)) {
		inPath = userHomeDir() + inPath[5:]
	}

	// 路径模板字符串替换
	inPath = os.ExpandEnv(inPath)

	// 判断路径是否是决定路径
	if filepath.IsAbs(inPath) {
		// 清理路径中的多余字符，比如 /// 或 ../ 或 ./
		return filepath.Clean(inPath)
	}

	// 转换为绝对路径
	p, err := filepath.Abs(inPath)
	if err == nil {
		return filepath.Clean(p)
	}

	return ""
}

type PathInfo struct {
	Name  string // 路径名称
	IsDir bool   // 是否是目录
}

// 获得最后一个字符.
func LastChar(str string) uint8 {
	if str == "" {
		panic("The length of the string can't be 0")
	}
	return str[len(str)-1]
}

// JoinPaths 路径合并.
func JoinPaths(absolutePath, relativePath string) string {
	if relativePath == "" {
		return absolutePath
	}

	finalPath := path.Join(absolutePath, relativePath)
	if LastChar(relativePath) == '/' && LastChar(finalPath) != '/' {
		return finalPath + "/"
	}
	return finalPath
}

// PathExists 判断文件路径是否存在.
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetWd 获得应用程序当前路径.
func GetWd() string {
	wd, _ := os.Getwd()
	return wd
}

// GetHomeDir 获得当前用户的$HOME目录.
func GetHomeDir() string {
	home, _ := homedir.Dir()
	return home
}

// Convert path to normal paths
func SlashAndCleanPath(path string) string {
	return filepath.ToSlash(filepath.Clean(path))
}

// 格式化路径字符串为 unix 下的绝对路径
// cleanPath is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//  1. Replace multiple slashes with a single slash.
//  2. Eliminate each . path name element (the current directory).
//  3. Eliminate each inner .. path name element (the parent directory)
//     along with the non-.. element that precedes it.
//  4. Eliminate .. elements that begin a rooted path:
//     that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned.
func CleanPath(p string) string {
	const stackBufSize = 128
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	// Reasonably sized buffer on stack to avoid allocations in the common case.
	// If a larger buffer is required, it gets allocated dynamically.
	buf := make([]byte, 0, stackBufSize)

	n := len(p)

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	// 如果路径不是以 / 开头，则可能需要重新分配缓存，并将第一个字符设置为 /
	if p[0] != '/' {
		r = 0

		if n+1 > stackBufSize {
			buf = make([]byte, n+1)
		} else {
			buf = buf[:n+1]
		}
		buf[0] = '/'
	}

	// 判断结尾是否包含 / 字符
	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp calls).
	// loop has no expensive function calls (except 1x make)		// So in contrast to the path package this loop has no expensive function
	// calls (except make, if needed).
	// 遍历源路径字符串
	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		// 如果最后一个字符是 .
		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		// 如果是 ./
		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		// 如果是 ../
		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if len(buf) == 0 {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// Real path element.
			// Add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// Copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// Re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	// If the original string was not modified (or only shortened at the end),
	// return the respective substring of the original string.
	// Otherwise return a new string from the buffer.
	if len(buf) == 0 {
		return p[:w]
	}
	return string(buf[:w])
}


// bufApp 将字符串 s[:w]拷贝到缓存 buf 中，重新计算缓存的大小刚好能容纳字符串的长度
// 并修改指定索引位置的字符值
//
// buf 为指针的原因是这是一个可修改值，执行后会修改字节数组的内容
func bufApp(buf *[]byte, s string, w int, c byte) {
	b := *buf
	if len(b) == 0 {
		// No modification of the original string so far.
		// If the next character is the same as in the original string, we do
		// not yet have to allocate a buffer.
		if s[w] == c {
			return
		}

		// 重新计算缓存的大小
		// Otherwise use either the stack buffer, if it is large enough, or
		// allocate a new buffer on the heap, and copy all previous characters.
		length := len(s)
		if length > cap(b) {
			// 源字符串超过缓存 bug 的长度，则重新创建一个新的数组
			*buf = make([]byte, length)
		} else {
			// 否则减少缓存的大小
			*buf = (*buf)[:length]
		}
		b = *buf

		// 将字符串 [ 0, w) 拷贝到缓存中
		copy(b, s[:w])
	}

	// 并将 w 索引修改为指定字符 c
	b[w] = c
}
