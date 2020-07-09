package profefe

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestCMD_Execute(t *testing.T) {
	// get `go` executable path
	goExecutable, _ := exec.LookPath("go")
	// construct `go version` command
	cmdGoVer := &exec.Cmd{
		Path:   goExecutable,
		Args:   []string{goExecutable, "tool", "pprof", "-http=:8082", "http://localhost:8081/api/0/profiles/bs3lvvbc1osn063vpf80"},
		Stdout: os.Stdout,
		Stderr: os.Stdout,
	}
	// run `go version` command
	if err := cmdGoVer.Run(); err != nil {
		fmt.Println("Error:", err);
	}
}