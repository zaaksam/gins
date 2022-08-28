package cmdutil

import (
	"flag"
	"os/exec"
	"time"
)

// Exec 执行命令，默认 10 秒超时
func Exec(name string, args []string, timeout ...time.Duration) (out []byte, err error) {
	c := exec.Command(name, args...)

	to := 10 * time.Second
	if len(timeout) == 1 {
		to = timeout[0]
	}

	stopTimer := time.AfterFunc(to, func() {
		c.Process.Kill()
	})

	out, err = c.CombinedOutput()
	if err == flag.ErrHelp {
		err = nil
	}

	stopTimer.Stop()

	return
}
