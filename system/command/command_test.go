package command

import (
	"fmt"
	"os/exec"
	"runtime"
	"testing"

	"github.com/duke-git/lancet/v2/system"
	"github.com/stretchr/testify/assert"
)

func TestRunBashCommand(t *testing.T) {
	sysType := runtime.GOOS
	if sysType != "windows" {
		out, errout, err := RunBashCommand("echo 1")
		assert.Equal(t, nil, err)
		assert.Equal(t, "1\n", out)
		assert.Equal(t, "", errout)
	}
}

// TestExecCommand 执行shell命令，返回命令的标准输出和标准错误字符串，
// 如果发生错误则返回错误。
// 参数command是一个完整的命令字符串，例如：ls -a（Linux）、dir（Windows）、ping 127.0.0.1。
// 在Linux中，使用/bin/bash -c来执行命令；
// 在Windows中，使用powershell.exe来执行命令。
// 函数的第二个参数是cmd选项控制参数。类型为func(*exec.Cmd)。可以通过此参数设置cmd属性。
// type (
//
//	Option func(*exec.Cmd)
//
// )
// func ExecCommand(command string, opts ...Option) (stdout, stderr string, err error)
func TestExecCommand(t *testing.T) {
	// linux or mac
	stdout, stderr, err := system.ExecCommand("ls", func(cmd *exec.Cmd) {
		cmd.Dir = "/tmp"
	})
	fmt.Println("std out: ", stdout)
	fmt.Println("std err: ", stderr)
	assert.Equal(t, "", stderr)

	// windows
	stdout, stderr, err = system.ExecCommand("dir")
	fmt.Println("std out: ", stdout)
	fmt.Println("std err: ", stderr)

	// error command
	stdout, stderr, err = system.ExecCommand("abc")
	fmt.Println("std out: ", stdout)
	fmt.Println("std err: ", stderr)
	if err != nil {
		fmt.Println(err.Error())
	}
}
