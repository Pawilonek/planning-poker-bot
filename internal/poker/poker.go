package poker

import (
	"fmt"
	"github.com/Pawilonek/scrumpoke/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"strings"
)

const (
	gameNotExists   = "Planning do not exists. Create new one with /poker command."
	gameVotedFormat = "You have voted [%s] in planning [%s]"
)

type (
	Poker struct {
		bot         *tgbotapi.BotAPI
		jiraConfig  config.Jira
		activeGames games
	}
)

func NewPoker(b *tgbotapi.BotAPI, j config.Jira) (p Poker) {
	p = Poker{bot: b, jiraConfig: j, activeGames: make(games, 0)}
	return
}

func (p *Poker) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	for update := range p.bot.GetUpdatesChan(u) {
		if update.CallbackQuery != nil {
			p.handleCallbackData(update)
			continue
		}

		if update.Message != nil && update.Message.IsCommand() {
			p.handleBotCommand(update)
		}
	}
}

func (p *Poker) handleBotCommand(u tgbotapi.Update) {
	m := tgbotapi.NewMessage(u.Message.Chat.ID, "")

	if strings.HasPrefix(u.Message.Command(), "poker") {
		jiraIssue := p.parseJiraIssue(u.Message.Text)
		hash := newHashFromBotCommandMessage(u.Message, jiraIssue)
		currentGame := p.activeGames.get(hash)
		if currentGame == nil {
			currentGame = newGame(jiraIssue, u.Message)
			p.activeGames.add(hash, currentGame)
		}

		keyboard := newKeyboard(hash, currentGame)
		m.ReplyMarkup = keyboard.buildMarkup()
		m.Text = keyboard.buildText()
	}

	if _, err := p.bot.Send(m); err != nil {
		log.Println(err)
	}
}

func (p *Poker) handleCallbackData(u tgbotapi.Update) {
	callbackData := newCallBackDataFromString(u.CallbackQuery.Data)
	currentGame := p.activeGames.get(callbackData.hash)
	if currentGame == nil {
		callback := tgbotapi.NewCallback(u.CallbackQuery.ID, gameNotExists)
		if _, err := p.bot.Request(callback); err != nil {
			log.Println(err)
		}
		return
	}
	f := newVote(u.CallbackQuery.From, callbackData)
	currentGame.addVote(f)

	keyboard := newKeyboard(callbackData.hash, currentGame)

	m := tgbotapi.NewEditMessageTextAndMarkup(
		u.CallbackQuery.Message.Chat.ID,
		u.CallbackQuery.Message.MessageID,
		keyboard.buildText(),
		keyboard.buildMarkup(),
	)

	if _, err := p.bot.Send(m); err != nil {
		log.Println(err)
	}

	callback := tgbotapi.NewCallback(u.CallbackQuery.ID, fmt.Sprintf(gameVotedFormat, callbackData.data, callbackData.hash))
	if _, err := p.bot.Request(callback); err != nil {
		log.Println(err)
	}
}

func (p *Poker) parseJiraIssue(t string) string {
	return strings.ReplaceAll(strings.TrimPrefix(t, "/poker"), "_", "-")
}
