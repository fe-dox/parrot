package telegraph

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
)

type (
	CallbackID     int64
	CallbackItemID int64
	CallbackStack  map[CallbackID]Callback
	Callback       map[CallbackItemID]CallbackItem
	CallbackItem   struct {
		Command string
		Data    string
	}
)

type CallbackResolver interface {
	AddCallback() CallbackID
	ClearCallback(id CallbackID)
	GetItem(csid CallbackID, ciid CallbackItemID) CallbackItem
}

type CallbackItemsResolver interface {
	AddItem(i CallbackItem) CallbackItemID
}

func (c CallbackStack) AddCallback() CallbackID {
	randID := CallbackID(rand.Int63())
	c[randID] = make(Callback)
	return randID
}

func (c CallbackStack) ClearCallback(id CallbackID) {
	delete(c, id)
	return
}

func (c CallbackStack) GetItem(csid CallbackID, ciid CallbackItemID) CallbackItem {
	return c[csid][ciid]
}

func (c Callback) AddItem(i CallbackItem) CallbackItemID {
	randID := CallbackItemID(rand.Int63())
	c[randID] = i
	return randID
}

func PrepareString(csid CallbackID, ciid CallbackItemID) string {
	return fmt.Sprintf("%v-%v", csid, ciid)
}

func DecodeString(s string) (csid CallbackID, ciid CallbackItemID, err error) {
	ss := strings.SplitN(s, "-", 2)
	pint, err := strconv.ParseInt(ss[0], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	csid = CallbackID(pint)
	pint, err = strconv.ParseInt(ss[1], 10, 64)
	if err != nil {
		return 0, 0, err
	}
	ciid = CallbackItemID(pint)
	return
}

func splitCommandAndData(s string) (string, string) {
	i := strings.Index(s, "-")
	return s[:i], s[i+1:]
}
