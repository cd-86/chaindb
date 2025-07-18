// 一台机器上只允许有一个 ChainDB 进程.

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/gofrs/flock"
)

func assertNoOtherChainDBInstance() {
	lock_path := filepath.Join(
		func() string {
			switch OS := runtime.GOOS; OS {
			case "windows":
				return `C:\ProgramData`
			case "linux":
				return "/var/lock"
			default:
				panic(
					fmt.Sprintf(
						"平台 OS (%s) 上的 ChainDB 进程锁 (用于确保本机上只有一个 ChainDB 进程) 暂时没有得到支持",
						OS,
					),
				)
			}
		}(),
		"shynur", "app", "chaindb", "proc-mutex.lock",
	)
	if err := os.MkdirAll(filepath.Dir(lock_path), 0777); err != nil {
		panic(err)
	}

	ok, err := flock.New(lock_path).TryLock()
	if err != nil {
		lock_file, err := os.OpenFile(lock_path, os.O_CREATE, 0666)
		if err != nil {
			panic(err)
		}
		lock_file.Close()
		ok, err = flock.New(lock_path).TryLock()
	}

	if err != nil {
		panic(err)
	}
	if !ok {
		fmt.Fprintf(os.Stderr, "本机上有其它 ChainDB 进行在运行\n")
		os.Exit(1)
	}
}
