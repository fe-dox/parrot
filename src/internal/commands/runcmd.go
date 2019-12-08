package commands

import (
	"fmt"
	"os/exec"
	"strings"
	"syscall"
)

func StartCommand(command string) string {
	args := strings.Fields(command)
	if len(args) < 2 {
		args = append(args, "")
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Failed: %v", err)
	}
	return fmt.Sprintf("Combined out:\n%s\n", string(out))
}

func RunCommand(command string) string {
	fullCmd := append([]string{"/C"}, strings.Fields(command)...)
	cmd := exec.Command("cmd", fullCmd...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Sprintf("Failed: %v", err)
	}
	return fmt.Sprintf("Combined out:\n%s\n", string(out))
}
