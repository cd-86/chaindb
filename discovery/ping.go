package discovery

import (
	"fmt"
	"os/exec"
	"runtime"
)

func ping(ip string) (err error) {
	switch OS := runtime.GOOS; OS {
	case "windows":
		err = exec.Command("ping", "-n", "1", ip).Run()
	case "linux":
		err = exec.Command("ping", "-c", "1", ip).Run()
	default:
		panic(
			fmt.Sprintf(
				"平台 OS (%s) 上的 `ping' 暂时没有得到 ChainDB 的支持",
				OS,
			),
		)
	}

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return
		} else {
			panic("`exec ping' 失败")
		}
	}
	return
}
