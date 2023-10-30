package poker

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	callbackDataDelimiter = ":"
	callbackDataFormat    = "%s" + callbackDataDelimiter + "%s"
)

var (
	keys = [][]string{
		{"1", "2", "3"},
		{"5", "8", "13"},
		{"20", "40", "100"},
		{"coffee", "?"},
	}
)

type CallbackData struct {
	hash Hash
	data string
}

type Keyboard struct {
	rows [][]string
	hash Hash
}

func NewCallBackData(data string) CallbackData {
	x := strings.Split(data, callbackDataDelimiter)
	return CallbackData{
		hash: Hash(x[0]),
		data: x[1],
	}
}

func (c CallbackData) String() string {
	return fmt.Sprintf(callbackDataFormat, c.hash, c.data)
}

func (c CallbackData) GetData() string {
	return c.data
}

func (c CallbackData) GetHash() Hash {
	return c.hash
}

func NewKeyboard(h Hash) Keyboard {
	return Keyboard{
		hash: h,
		rows: keys,
	}
}

func (k Keyboard) BuildMarkup() tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, column := range k.rows {
		buttons := make([]tgbotapi.InlineKeyboardButton, 0)
		for _, v := range column {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v, k.createCallbackData(v)))
		}
		rows = append(rows, [][]tgbotapi.InlineKeyboardButton{buttons}...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

func (k Keyboard) createCallbackData(data string) string {
	return fmt.Sprintf(callbackDataFormat, k.hash, data)
}
