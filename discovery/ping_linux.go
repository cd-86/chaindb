package discovery

import "os/exec"

func ping(ip string) (err error) {
	err = exec.Command("ping", "-c", "1", ip).Run()

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return
		} else {
			panic("`exec ping' 失败")
		}
	}

	return
}
