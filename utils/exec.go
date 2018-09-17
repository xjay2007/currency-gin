package utils

import (
	"strings"
	"os/exec"
)

func ExecCmd(name string, arg ...string) (bool, string) {
	Info("execCmd:", name, strings.Join(arg, " "))
	cmd := exec.Command(name, arg...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)
	if err != nil {
		Info(outStr)
		Error(err)
		return false, outStr
	}
	//Info(outStr)
	return true, outStr
}