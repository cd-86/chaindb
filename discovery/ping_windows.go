package discovery

import "os/exec"

func ping(ip string) (err error) {
	proc_ping := exec.Command("ping", "-n", "1", ip)
	proc_ping.SysProcAttr.HideWindow = true
	err = proc_ping.Run()

	if err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			return
		} else {
			panic("`exec ping' 失败")
		}
	}

	return
}
