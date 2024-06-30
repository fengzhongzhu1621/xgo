package command

import (
	"bufio"
	"bytes"
	"io"
	"os"
	"os/exec"
)

// 阻塞运行 linux bash 命令.
func RunBashCommand(command string) (string, string, error) {
	// 返回一个*Cmd，用于使用给出的参数执行name指定的程序
	cmd := exec.Command("/bin/bash", "-c", command)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	// 阻塞执行命令
	err := cmd.Run()
	// 返回执行结果
	return stdout.String(), stderr.String(), err
}

// 后台运行linux bash命令.
func RunBashCommandBackground(commandName string) *exec.Cmd {
	cmd := exec.Command("/bin/bash", "-c", commandName)

	go func() {
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
	}()

	return cmd
}

// 直接在当前目录使用并返回结果.
func CmdInDir(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	// 显示运行的命令
	// fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	err := cmd.Start()
	if err != nil {
		return "", err
	}
	// 阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err = cmd.Wait()

	// 返回执行结果
	return out.String(), err
}

// 在命令位置使用并返回结果.
func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	// fmt.Println("CmdAndChangeDir", dir, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir

	err := cmd.Start()
	if err != nil {
		return "", err
	}
	// 阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err = cmd.Wait()
	// 返回执行结果
	return out.String(), err
}

// 在命令位置使用并实时输出每行结果.
func CmdAndChangeDirToShow(dir string, commandName string, params []string) error {
	cmd := exec.Command(commandName, params...)
	// fmt.Println("CmdAndChangeDirToFile", dir, cmd.Args)
	// StdoutPipe 方法返回一个在命令Start后与命令标准输出关联的管道。
	// Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		// fmt.Println("cmd.StdoutPipe: ", err)
		return err
	}
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err
	}

	// 创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	// 实时循环读取输出流中的一行内容
	for {
		_, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		// fmt.Println(line)
	}
	// 阻塞直到该命令执行完成，该命令必须是被Start方法开始执行的
	err = cmd.Wait()

	return err
}

// 命令加载进程类.
type LaunchedProcess struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

// 加载一个命令.
func LaunchCmd(commandName string, commandArgs []string, env []string) (*LaunchedProcess, error) {
	cmd := exec.Command(commandName, commandArgs...)
	cmd.Env = env

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, err
	}

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	// 执行命令，需要等待命令执行完成
	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	return &LaunchedProcess{cmd, stdin, stdout, stderr}, err
}

// 将多个参数添加到参数数组中.
func AppendArgs(dst, src []interface{}) []interface{} {
	if len(src) == 1 {
		return AppendArg(dst, src[0])
	}

	dst = append(dst, src...)
	return dst
}

// 将单个参数 arg 添加到参数数组 dst 中.
func AppendArg(dst []interface{}, arg interface{}) []interface{} {
	switch arg := arg.(type) {
	case []string:
		for _, s := range arg {
			dst = append(dst, s)
		}
		return dst
	case []interface{}:
		dst = append(dst, arg...)
		return dst
	case map[string]interface{}:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	case map[string]string:
		for k, v := range arg {
			dst = append(dst, k, v)
		}
		return dst
	default:
		return append(dst, arg)
	}
}
