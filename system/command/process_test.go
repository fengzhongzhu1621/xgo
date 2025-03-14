package command

import (
	"fmt"
	"testing"
	"time"

	"github.com/duke-git/lancet/v2/system"
)

// StartProcess Start a new process with the specified name and arguments.
// func StartProcess(command string, args ...string) (int, error)
// func StopProcess(pid int) error
// func KillProcess(pid int) error
func TestStartProcess(t *testing.T) {
	pid, err := system.StartProcess("sleep", "2")
	if err != nil {
		return
	}

	fmt.Println(pid)
	time.Sleep(1 * time.Second)

	err = system.StopProcess(pid)
	fmt.Println(err)

	err = system.KillProcess(pid)
	fmt.Println(err)
}

// GetProcessInfo Retrieves detailed process information by pid.
// func GetProcessInfo(pid int) (*ProcessInfo, error)
func TestGetProcessInfo(t *testing.T) {
	pid, err := system.StartProcess("ls", "-l")
	if err != nil {
		return
	}
	fmt.Println(pid)

	processInfo, err := system.GetProcessInfo(pid)
	if err != nil {
		return
	}

	fmt.Println(processInfo)
}
