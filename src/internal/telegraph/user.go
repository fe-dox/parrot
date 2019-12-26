package telegraph

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io/ioutil"
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
	files, err := ioutil.ReadDir(dir.path)
	for _, info := range files {
		if info.IsDir() {
			name := info.Name()
			dir.innerDirs = append(dir.innerDirs, Item{
				name: name,
				path: filepath.Clean(dir.path + "\\" + name),
				info: info,
			})
		} else {
			dir.innerFiles = append(dir.innerFiles, Item{
				path: filepath.Clean(dir.path + "\\" + info.Name()),
				info: info,
			})
		}

	}

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

func (u *User) ScanCurrentPath() (Directory, error) {
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

func (u *User) SetPath(s string) error {
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

	parentDir := filepath.Clean(d.path + "\\..")

	functionalRow := make([]tgbotapi.InlineKeyboardButton, 2)
	functionalRow[0] = t.callbackStack.CreateButton(cbID, "🔝", FilesystemPathRequest, parentDir)
	functionalRow[1] = t.callbackStack.CreateButton(cbID, "↩", FilesystemPathRequest, d.path)

	chunkedInnerDirs := chunkArray(d.innerDirs, 5)

	allRows := make([][]tgbotapi.InlineKeyboardButton, len(chunkedInnerDirs)+1)
	allRows[0] = functionalRow

	for j, x := range chunkedInnerDirs {
		dataRow := make([]tgbotapi.InlineKeyboardButton, len(x))
		for i, v := range x {
			dataRow[i] = t.callbackStack.CreateButton(cbID, v.name, FilesystemPathRequest, v.path)
		}
		allRows[j+1] = dataRow
	}

	return tgbotapi.InlineKeyboardMarkup{InlineKeyboard: allRows}
}

func chunkArray(arr []Item, chunkSize int) [][]Item {
	var divided [][]Item
	for i := 0; i < len(arr); i += chunkSize {
		end := i + chunkSize
		if end > len(arr) {
			end = len(arr)
		}
		divided = append(divided, arr[i:end])
	}
	return divided
}
