package fsnotify

import (
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/fsnotify/fsnotify"
)

func TestMonitorDir(t *testing.T) {
	// 创建一个新的 fsnotify.Watcher 实例。
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// 添加当前目录 (.) 到监控列表中。
	dir := "."
	err = watcher.Add(dir)
	if err != nil {
		log.Fatal(err)
	}

	// 启动一个 goroutine 来监听和处理来自 watcher.Events 和 watcher.Errors 的事件。
	go func() {
		for {
			select {
			case event := <-watcher.Events:
				fmt.Println("Event:", event)
				// 检查文件/目录创建事件
				if event.Op&fsnotify.Create == fsnotify.Create {
					fmt.Println("File created:", event.Name)
				}
				// 检查文件/目录删除事件
				if event.Op&fsnotify.Remove == fsnotify.Remove {
					fmt.Println("File removed:", event.Name)
				}
				// 检查文件修改事件
				if event.Op&fsnotify.Write == fsnotify.Write {
					fmt.Println("File modified:", event.Name)
				}
			case err := <-watcher.Errors:
				fmt.Println("Error:", err)
			}
		}
	}()

	// 模拟一些文件操作
	// 程序在启动后等待 2 秒，创建一个名为 test.txt 的文件。
	time.Sleep(2 * time.Second)
	fmt.Println("Creating file test.txt...")
	_, err = createFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
	writeToFile("test.txt")

	// 再等待 2 秒，删除 test.txt 文件。
	time.Sleep(2 * time.Second)
	fmt.Println("Deleting file test.txt...")
	err = deleteFile("test.txt")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(2 * time.Second)
}

func createFile(filename string) (*os.File, error) {
	return os.Create(filename)
}

func deleteFile(filename string) error {
	return os.Remove(filename)
}

func writeToFile(filename string) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_WRONLY, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	file.WriteString("Hello, fsnotify!\n")
}
