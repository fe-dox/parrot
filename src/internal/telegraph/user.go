package telegraph

import (
	"fmt"
	"os"
	"path"
)

type User struct {
	authenticated bool
	currentPath   string
}

type Item struct {
	path string
	info os.FileInfo
}

type Directory struct {
	path       string
	info       os.FileInfo
	innerFiles []Item
	innerDirs  []Directory
}

type Actions interface {
	SetPath(s string) error
	ScanPath(s string) Directory
	ScanCurrentPath() Directory
}

func (u User) ScanPath(s string) Directory {

}

func (u User) ScanCurrentPath() Directory {
	return u.ScanPath(u.currentPath)
}

func (u User) SetPath(s string) error {
	if path.IsAbs(s) {
		s = path.Clean(s)
		pathStat, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("directory path %s doesn't exist", s)
			} else {
				return err
			}
		}
		if !pathStat.IsDir() {
			return fmt.Errorf("path %v is not directory", s)
		}
		u.currentPath = s
		return nil
	} else {
		s = path.Join(u.currentPath, s)
		s = path.Clean(s)
		pathStat, err := os.Stat(s)
		if err != nil {
			if os.IsNotExist(err) {
				return fmt.Errorf("directory path %s doesn't exist", s)
			} else {
				return err
			}
		}
		if !pathStat.IsDir() {
			return fmt.Errorf("path %v is not directory", s)
		}
		u.currentPath = s
		return nil
	}
}
