package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

func listAllDrives() ([]string, error) {
	kernel32, err := syscall.LoadDLL("kernel32.dll")
	if err != nil {
		return []string{}, err
	}
	getLogicalDriveStringsHandle, err := kernel32.FindProc("GetLogicalDriveStringsA")
	if err != nil {
		return []string{}, err
	}
	buffer := [1024]byte{}
	bufferSize := uint32(len(buffer))
	drives := []string{}

	hr, _, _ := getLogicalDriveStringsHandle.Call(uintptr(unsafe.Pointer(&bufferSize)), uintptr(unsafe.Pointer(&buffer)))
	if hr == 0 {
		return drives, fmt.Errorf("an error occured")
	} else {

		bytes.Split(buffer, []byte{58, 92, 0})
		return drives, nil
	}
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

func splitPath(path string) (dir, file string) {
	i := strings.LastIndex(path, "\\")
	return path[:i+1], path[i+1:]
}
