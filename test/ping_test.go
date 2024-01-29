package test

import (
	"bytes"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

func TestPingUse(t *testing.T) {
	use, i := PingUse("82.157.208.214", t)
	fmt.Println((i))
	fmt.Println(use)
}

func PingUse(ip string, t *testing.T) (bool, []byte) {
	var command *exec.Cmd
	switch runtime.GOOS {
	case "windows":
		command = exec.Command("cmd", "/c", "ping -n 1 -w 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	case "darwin":
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -W 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	default: //linux
		command = exec.Command("/bin/bash", "-c", "ping -c 1 -w 1 "+ip+" && echo true || echo false") //ping -c 1 -i 0.5 -t 4 -W 2 -w 5 "+ip+" >/dev/null && echo true || echo false"
	}
	outInfo := bytes.Buffer{}
	command.Stdout = &outInfo
	err := command.Start()
	if err != nil {
		return false, nil
	}
	if err = command.Wait(); err != nil {
		return false, nil
	} else {
		if strings.Contains(outInfo.String(), "true") && strings.Count(outInfo.String(), ip) > 2 {
			return true, outInfo.Bytes()
		} else {
			return false, nil
		}
	}
}
