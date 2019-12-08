package commands

import (
	"bytes"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/kbinani/screenshot"
	"image/png"
)

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
