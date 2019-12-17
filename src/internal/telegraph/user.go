package telegraph

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
)

type User struct {
	authenticated bool
	currentPath   string
	currentDir    Directory
}

type Item struct {
	name string
	path string
	info os.FileInfo
}

type Directory struct {
	path       string
	info       os.FileInfo
	innerFiles []Item
	innerDirs  []Item
}

type Actions interface {
	SetPath(s string) error
	ScanPath(s string) Directory
	ScanCurrentPath() Directory
}

func (u User) ScanPath(s string) (Directory, error) {
	dir := Directory{
		path:       s,
		info:       nil,
		innerFiles: nil,
		innerDirs:  nil,
	}
	fileInfo, err := os.Stat(s)
	dir.info = fileInfo
	if err != nil {
		return Directory{
			path:       "",
			info:       nil,
			innerFiles: nil,
			innerDirs:  nil,
		}, err
	}
	err = filepath.Walk(s, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			name := filepath.FromSlash(path)
			dir.innerDirs = append(dir.innerDirs, Item{
				name: name,
				path: path,
				info: info,
			})
		} else {
			dir.innerFiles = append(dir.innerFiles, Item{
				path: path,
				info: info,
			})
		}
		return err
	})

	if err != nil {
		return Directory{
			path:       "",
			info:       nil,
			innerFiles: nil,
			innerDirs:  nil,
		}, err
	}
	return dir, nil
}

func (u User) ScanCurrentPath() (Directory, error) {
	tmpDir, err := u.ScanPath(u.currentPath)
	if err != nil {
		return Directory{
			path:       "",
			info:       nil,
			innerFiles: nil,
			innerDirs:  nil,
		}, err
	}
	u.currentDir = tmpDir
	return u.currentDir, nil
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
