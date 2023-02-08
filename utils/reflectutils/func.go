package reflectutils

import (
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

// Func contains runtime information about a function.
type Func struct {
	// Name of the function.
	Name string

	// Name of the package in which this function is defined.
	Package string

	// Path to the file in which this function is defined.
	File string

	// Line number in the file at which this function is defined.
	Line int
}

// String returns a string representation of the function.
func (f *Func) String() string {
	return fmt.Sprint(f)
}

// Format implements fmt.Formatter for Func, printing a single-line
// representation for %v and a multi-line one for %+v.
// 设置Sprint的格式
func (f *Func) Format(w fmt.State, c rune) {
	// %+v     打印结构体时，会添加字段名     Printf("%+v", people)  {Name:zhangsan}
	if w.Flag('+') && c == 'v' {
		// "path/to/package".MyFunction
		// 	path/to/file.go:42
		// %q      单引号围绕的字符字面值，由Go语法安全地转义 Printf("%q", 0x4E2D)        '中'
		// %v      相应值的默认格式。            Printf("%v", people)   {zhangsan}，
		fmt.Fprintf(w, "%q.%v", f.Package, f.Name)
		fmt.Fprintf(w, "\n\t%v:%v", f.File, f.Line)
	} else {
		// "path/to/package".MyFunction (path/to/file.go:42)
		fmt.Fprintf(w, "%q.%v (%v:%v)", f.Package, f.Name, f.File, f.Line)
	}
}

// InspectFunc inspects and returns runtime information about the given
// function.
// 获得函数的运行时信息
func InspectFunc(function interface{}) *Func {
	fptr := reflect.ValueOf(function).Pointer()
	return InspectFuncPC(fptr)
}

// InspectFuncPC inspects and returns runtime information about the function
// at the given program counter address.
func InspectFuncPC(pc uintptr) *Func {
	// 返回一个表示调用栈标识符pc对应的调用栈的*Func;如果该调用栈标识符没有对应的调用栈，函数会返回nil。
	// 通过reflect的ValueOf().Pointer作为入参，获取函数地址、文件行、函数名等信息
	f := runtime.FuncForPC(pc)
	if f == nil {
		return nil
	}
	// 获得函数所在的包名和函数名
	pkgName, funcName := splitFuncName(f.Name())
	// 获得函数所在的文件名和行号
	fileName, lineNum := f.FileLine(pc)
	return &Func{
		Name:    funcName,
		Package: pkgName,
		File:    fileName,
		Line:    lineNum,
	}
}

const _vendor = "/vendor/"

// 将函数的运行时name分割为包名和函数名
func splitFuncName(function string) (pname string, fname string) {
	if len(function) == 0 {
		return
	}

	// We have something like "path.to/my/pkg.MyFunction". If the function is
	// a closure, it is something like, "path.to/my/pkg.MyFunction.func1".

	idx := 0

	// Everything up to the first "." after the last "/" is the package name.
	// Everything after the "." is the full function name.
	if i := strings.LastIndex(function, "/"); i >= 0 {
		idx = i
	}
	if i := strings.Index(function[idx:], "."); i >= 0 {
		idx += i
	}
	pname, fname = function[:idx], function[idx+1:]

	// The package may be vendored.
	if i := strings.Index(pname, _vendor); i > 0 {
		pname = pname[i+len(_vendor):]
	}

	// Package names are URL-encoded to avoid ambiguity in the case where the
	// package name contains ".git". Otherwise, "foo/bar.git.MyFunction" would
	// mean that "git" is the top-level function and "MyFunction" is embedded
	// inside it.
	// 将QueryEscape转码的字符串还原。它会把%AB改为字节0xAB，将’+‘改为’ '。
	// 如果有某个%后面未跟两个十六进制数字，会返回错误。
	if unescaped, err := url.QueryUnescape(pname); err == nil {
		pname = unescaped
	}

	return
}

// MakeAddressable returns a value that is always addressable.
// It returns the input verbatim if it is already addressable,
// otherwise it creates a new value and returns an addressable copy.
func MakeAddressable(v reflect.Value) reflect.Value {
	if v.CanAddr() {
		return v
	}
	vc := reflect.New(v.Type()).Elem()
	vc.Set(v)
	return vc
}
