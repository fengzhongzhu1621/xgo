package command

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
)

/**
 * 运行linux bash命令
 */
func RunBashCommand(command string) (error, string, string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	return err, stdout.String(), stderr.String()
}


// 后台运行linux bash命令
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

/**
 * 直接在当前目录使用并返回结果
 */
func Cmd(commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("Cmd", cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}


/**
 * 在命令位置使用并返回结果
 */
func CmdAndChangeDir(dir string, commandName string, params []string) (string, error) {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDir", dir, cmd.Args)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err := cmd.Start()
	if err != nil {
		return "", err
	}
	err = cmd.Wait()
	return out.String(), err
}

/**
 * 在命令位置使用并实时输出每行结果
 */
func CmdAndChangeDirToShow(dir string, commandName string, params []string) error {
	cmd := exec.Command(commandName, params...)
	fmt.Println("CmdAndChangeDirToFile", dir, cmd.Args)
	// StdoutPipe方法返回一个在命令Start后与命令标准输出关联的管道。
	// Wait方法获知命令结束后会关闭这个管道，一般不需要显式的关闭该管道。
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println("cmd.StdoutPipe: ", err)
		return err
	}
	cmd.Stderr = os.Stderr
	cmd.Dir = dir
	err = cmd.Start()
	if err != nil {
		return err
	}
	//创建一个流来读取管道内内容，这里逻辑是通过一行一行的读取的
	reader := bufio.NewReader(stdout)
	//实时循环读取输出流中的一行内容
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		fmt.Println(line)
	}
	err = cmd.Wait()
	return err
}


// 命令加载进程类
type LaunchedProcess struct {
	Cmd    *exec.Cmd
	Stdin  io.WriteCloser
	Stdout io.ReadCloser
	Stderr io.ReadCloser
}

/**
 * 加载一个命令
 */
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
