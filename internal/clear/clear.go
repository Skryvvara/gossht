package clear

import (
	"os"
	"os/exec"
	"runtime"
)

func isWindows() bool {
  return runtime.GOOS == "windows"
}

func CallClear() {
  var cmd *exec.Cmd
  if isWindows() {
    cmd = exec.Command("cmd", "/c", "cls")
  } else {
    cmd = exec.Command("clear")
  }

  cmd.Stdout = os.Stdout
  cmd.Run()
}
