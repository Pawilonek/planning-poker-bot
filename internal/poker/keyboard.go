package poker

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	callbackDataDelimiter = ":"
	callbackDataFormat    = "%s" + callbackDataDelimiter + "%s"

	planningDefaultMessageFormat      = "Planning ID: %s"
	planningWithJiraTaskMessageFormat = "Planning: %s"
	planningCurrentVotersFormat       = "Voters:"
)

var (
	keys = [][]string{
		{"1", "2", "3"},
		{"5", "8", "13"},
		{"20", "40", "100"},
		{"coffee", "?"},
	}
)

type (
	callbackData struct {
		hash gameHash
		data string
	}
	keyboard struct {
		rows     [][]string
		gameHash gameHash
		game     game
	}
)

func newCallBackDataFromString(data string) callbackData {
	x := strings.Split(data, callbackDataDelimiter)
	return callbackData{
		hash: gameHash(x[0]),
		data: x[1],
	}
}

func newCallbackData(g gameHash, data string) callbackData {
	return callbackData{
		hash: g,
		data: data,
	}
}

func (c callbackData) string() string {
	return fmt.Sprintf(callbackDataFormat, c.hash, c.data)
}

func newKeyboard(h gameHash, g *game) keyboard {
	return keyboard{
		rows:     keys,
		gameHash: h,
		game:     *g,
	}
}

func (k keyboard) buildText() (text string) {
	text = fmt.Sprintf(planningDefaultMessageFormat, k.gameHash)
	if k.game.jiraIssue != "" {
		text = fmt.Sprintf(planningWithJiraTaskMessageFormat, k.game.jiraIssue)
	}

	voters := k.game.listVoters()
	if len(voters) > 0 {
		text = text + "\n" + planningCurrentVotersFormat
		for _, v := range voters {
			text = text + fmt.Sprintf("\n  - %s", v.UserName)
		}
	}

	return
}

func (k keyboard) buildMarkup() tgbotapi.InlineKeyboardMarkup {
	rows := make([][]tgbotapi.InlineKeyboardButton, 0)
	for _, column := range k.rows {
		buttons := make([]tgbotapi.InlineKeyboardButton, 0)
		for _, v := range column {
			buttons = append(buttons, tgbotapi.NewInlineKeyboardButtonData(v, newCallbackData(k.gameHash, v).string()))
		}
		rows = append(rows, [][]tgbotapi.InlineKeyboardButton{buttons}...)
	}
	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
