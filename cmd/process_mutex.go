// 一台机器上只允许有一个 ChainDB 进程.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func assertNoOtherChainDBInstance() {
	pid_file_path := filepath.Join(
		func() string {
			switch OS := runtime.GOOS; OS {
			case "windows":
				return `C:\ProgramData`
			case "linux":
				return "/tmp"
			default:
				panic(
					fmt.Sprintf(
						"平台 OS (%s) 上的 ChainDB 进程锁 (用于确保本机上只有一个 ChainDB 进程) 暂时没有得到支持",
						OS,
					),
				)
			}
		}(),
		"shynur", "app", "chaindb", "proc-mutex",
		fmt.Sprintf("%d.pid", os.Getpid()),
	)
	if err := os.MkdirAll(filepath.Dir(pid_file_path), 0777); err != nil {
		panic(err)
	}
	pid_file, err := os.Create(pid_file_path)
	if err != nil {
		panic(err)
	}
	pid_file.Close()

	pid_files, err := os.ReadDir(filepath.Dir(pid_file_path))
	if err != nil {
		panic(err)
	}
	for _, other_pid_file := range pid_files {
		pid, _ := strconv.Atoi(strings.TrimSuffix(other_pid_file.Name(), ".pid"))
		fmt.Println(pid)
	}

	if true {
		fmt.Fprintf(os.Stderr, "本机上有其它 ChainDB 进行在运行\n")
		os.Exit(1)
	}
}
