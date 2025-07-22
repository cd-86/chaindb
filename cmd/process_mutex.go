// 一台机器上只允许有一个 ChainDB 进程.

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	psutil_proc "github.com/shirou/gopsutil/v4/process"
)

var _ = func() any {
	assertNoOtherChainDBInstance()
	return nil
}()

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
	func() {
		pid_file, err := os.Create(pid_file_path)
		if err != nil {
			panic(err)
		} else {
			defer pid_file.Close()
		}
		ps_proc, err := psutil_proc.NewProcess(int32(os.Getpid()))
		if err != nil {
			panic(err)
		}
		cmdline, err := ps_proc.Cmdline()
		if err != nil {
			panic(err)
		}
		pid_file.WriteString(cmdline)
	}()

	pid_files, err := os.ReadDir(filepath.Dir(pid_file_path))
	if err != nil {
		panic(err)
	}
	for _, other_pid_file := range pid_files {
		if !strings.HasSuffix(other_pid_file.Name(), ".pid") {
			continue
		}
		pid, _ := strconv.Atoi(strings.TrimSuffix(other_pid_file.Name(), ".pid"))
		if pid == os.Getpid() {
			continue
		}

		other_pid_file_path := filepath.Join(filepath.Dir(pid_file_path), other_pid_file.Name())

		other_proc, err := psutil_proc.NewProcess(int32(pid))
		if err != nil {
			os.Remove(other_pid_file_path)
			continue
		}

		cmdline_in_file, err := os.ReadFile(other_pid_file_path)
		if err != nil {
			panic(err)
		}
		cmdline_of_proc, err := other_proc.Cmdline()
		if err != nil {
			panic(err)
		}
		if string(cmdline_in_file) != cmdline_of_proc {
			os.Remove(other_pid_file_path)
			continue
		}

		log.Fatalln("本机上有其它 ChainDB 进行在运行")
	}
}
