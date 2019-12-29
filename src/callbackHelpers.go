package main

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"math/rand"
	"strconv"
	"strings"
)

type (
	CallbackStackItemID int64
	CallbackItemID      int64
	CallbackStack       map[CallbackStackItemID]Callback
	Callback            map[CallbackItemID]CallbackItem
	CallbackItem        struct {
		Command string
		Data    string
	}
)

type CallbackResolver interface {
	AddCallback() CallbackStackItemID
	ClearCallback(id CallbackStackItemID)
	GetItem(csid CallbackStackItemID, ciid CallbackItemID) CallbackItem
	CreateButton(csid CallbackStackItemID, text string, command string, data string) tgbotapi.InlineKeyboardButton
	DecodeCallbackRequest(s string) (command, data string, err error)
}

type CallbackItemsResolver interface {
	AddItem(i CallbackItem) CallbackItemID
}

func (c CallbackStack) AddCallback() CallbackStackItemID {
	randID := CallbackStackItemID(rand.Int63())
	c[randID] = make(Callback)
	return randID
}

func (c CallbackStack) ClearCallback(id CallbackStackItemID) {
	delete(c, id)
	return
}

func (c CallbackStack) GetItem(csid CallbackStackItemID, ciid CallbackItemID) CallbackItem {
	return c[csid][ciid]
}

func (c Callback) AddItem(i CallbackItem) CallbackItemID {
	randID := CallbackItemID(rand.Int63())
	c[randID] = i
	return randID
}

func (c CallbackStack) CreateButton(csid CallbackStackItemID, text string, command string, data string) tgbotapi.InlineKeyboardButton {
	itemID := c[csid].AddItem(CallbackItem{
		Command: command,
		Data:    data,
	})
	return tgbotapi.NewInlineKeyboardButtonData(text, prepareString(csid, itemID))
}

func (c CallbackStack) DecodeCallbackRequest(s string) (command, data string, err error) {
	csid, ciid, err := DecodeString(s)
	if err != nil {
		return "", "", err
	}
	item := c.getCallbackItem(csid, ciid)
	if item.Data == "" && item.Command == "" {
		return "", "", fmt.Errorf("callback %v does not exist in current callback stack", csid)
	}
	return item.Command, item.Data, nil
}

func prepareString(csid CallbackStackItemID, ciid CallbackItemID) string {
	return fmt.Sprintf("%v-%v", csid, ciid)
}

func DecodeString(s string) (csid CallbackStackItemID, ciid CallbackItemID, err error) {
	ss := strings.SplitN(s, "-", 2)
	pint, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	csid = CallbackStackItemID(pint)
	pint, err = strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	ciid = CallbackItemID(pint)
	return
}

func (c CallbackStack) getCallbackItem(csid CallbackStackItemID, ciid CallbackItemID) CallbackItem {
	return c[csid][ciid]
}
