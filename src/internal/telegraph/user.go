package telegraph

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"os"
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
	ScanPath(s string) (Directory, error)
	ScanCurrentPath() (Directory, error)
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
			name := filepath.Base(path)
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
	if filepath.IsAbs(s) {
		s = filepath.Clean(s)
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
		s = filepath.Join(u.currentPath, s)
		s = filepath.Clean(s)
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

func (t Telegraphist) PrepareFilesystemKeyboard(d Directory) tgbotapi.InlineKeyboardMarkup {
	cbID := t.callbackStack.AddCallback()
	innerDirsLen := len(d.innerDirs)
	keyboardRow := make([]tgbotapi.InlineKeyboardButton, innerDirsLen+1)
	parentDir := filepath.Clean(d.path + "\\..")
	keyboardRow[0] = t.callbackStack.CreateButton(cbID, "..", FilesystemPathRequest, parentDir)
	for i, v := range d.innerDirs {
		if i+1 == innerDirsLen {
			v.name = "Rescan"
		}
		keyboardRow[i+1] = t.callbackStack.CreateButton(cbID, v.name, FilesystemPathRequest, d.path)
	}
	return tgbotapi.NewInlineKeyboardMarkup(keyboardRow)
}
