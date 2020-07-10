package port

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"syscall"
)

func  KillPort(port string) {

	if port == "" {
		port = "8082"
	}

	if _, err := strconv.Atoi(port); err != nil {
		os.Stderr.WriteString(fmt.Sprintf("Error: port argument is not a number.\n"))
		panic(err)
	}

	if runtime.GOOS == "windows" {
		command := fmt.Sprintf("(Get-NetTCPConnection -LocalPort %s).OwningProcess -Force", port)
		exec_cmd(exec.Command("Stop-Process", "-Id", command))
	} else {
		command := fmt.Sprintf("lsof -i tcp:%s | grep LISTEN | awk '{print $2}' | xargs kill -9", port)
		exec_cmd(exec.Command("bash", "-c", command))
	}
}
func exec_cmd(cmd *exec.Cmd) {
	var waitStatus syscall.WaitStatus
	if err := cmd.Run(); err != nil {
		if err != nil {
			os.Stderr.WriteString(fmt.Sprintf("Error: %s\n", err.Error()))
		}
		if exitError, ok := err.(*exec.ExitError); ok {
			waitStatus = exitError.Sys().(syscall.WaitStatus)
			fmt.Printf("Error during killing (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
		}
	} else {
		waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
		fmt.Printf("Port successfully killed (exit code: %s)\n", []byte(fmt.Sprintf("%d", waitStatus.ExitStatus())))
	}
}