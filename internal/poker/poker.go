package poker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/Pawilonek/scrumpoke/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"strings"
)

const (
	planningDefaultMessageFormat        = "Planning ID: %s:"
	planningWithJiraTaskMessageTemplate = "Currently planning: %s/%s:"
)

type (
	Hash string

	Poker struct {
		bot             *tgbotapi.BotAPI
		jiraConfig      config.Jira
		activePlannings plannings
	}

	plannings map[Hash]planning

	voteDetails struct {
		value     string
		createdAt int
		updatedAt int
	}

	vote struct {
		tgbotapi.User
		voteDetails
	}

	planning struct {
		votes     []vote
		createdAt int
		updatedAt int
		creator   *tgbotapi.User
		jiraIssue string
	}
)

func NewPoker(b *tgbotapi.BotAPI, j config.Jira) (p Poker) {
	activePlannings := make(plannings, 0)
	p = Poker{bot: b, jiraConfig: j, activePlannings: activePlannings}
	return
}

func (p *Poker) Run() {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updatesChannel := p.bot.GetUpdatesChan(u)

	for update := range updatesChannel {
		if update.CallbackQuery != nil {
			callbackData := NewCallBackData(update.CallbackQuery.Data)
			// Respond to the callback query, telling Telegram to show the user
			// a message with the data received.
			callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "You have voted ["+callbackData.GetData()+"] in planning ["+string(callbackData.GetHash())+"]")
			if _, err := p.bot.Request(callback); err != nil {
				panic(err)
			}

			// And finally, send a message containing the data received.
			msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf(
				"UserName: %s\nUserId: %d\nData: %s\nRespondsToMessage: %d\nChat: %s",
				update.CallbackQuery.From.String(),
				update.CallbackQuery.From.ID,
				update.CallbackQuery.Data,
				update.CallbackQuery.Message.MessageID,
				update.CallbackQuery.ChatInstance,
			))
			if _, err := p.bot.Send(msg); err != nil {
				panic(err)
			}

			continue
		}

		if update.Message == nil {
			fmt.Println("not a message")
			continue
		}

		if !update.Message.IsCommand() {
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")
		// Extract the command from the Message.
		if strings.HasPrefix(update.Message.Command(), "poker") {
			jiraIssue := p.getJiraIssue(update.Message.Text)
			hash := newHashFromBotCommandMessage(update.Message, jiraIssue)
			activePlanning := p.activePlannings.get(hash)
			if activePlanning == nil {
				activePlanning = newPlanning(jiraIssue, update.Message)
				p.activePlannings.add(hash, activePlanning)
			}
			msg.Text = fmt.Sprintf(planningDefaultMessageFormat, hash)
			if p.jiraConfig.TaskBrowseURL != "" && activePlanning.jiraIssue != "" {
				msg.Text = fmt.Sprintf(planningWithJiraTaskMessageTemplate, p.jiraConfig.TaskBrowseURL, activePlanning.jiraIssue)
			}
			keyboard := NewKeyboard(hash)
			msg.ReplyMarkup = keyboard.BuildMarkup()
		}

		if _, err := p.bot.Send(msg); err != nil {
			fmt.Println(err)
		}

	}
}

func (p *Poker) getJiraIssue(t string) string {
	return strings.ReplaceAll(strings.TrimPrefix(t, "/poker"), "_", "-")
}

func newPlanning(j string, m *tgbotapi.Message) *planning {
	return &planning{createdAt: m.Date, updatedAt: m.Date, creator: m.From, jiraIssue: j}
}

func newHashFromBotCommandMessage(m *tgbotapi.Message, jiraIssue string) (h Hash) {
	if m == nil {
		return ""
	}
	hash := md5.New()
	hash.Write([]byte(fmt.Sprintf(
		"%d%d%s",
		m.From.ID,
		m.Chat.ID,
		jiraIssue,
	)))
	return Hash(hex.EncodeToString(hash.Sum(nil)))
}

func (ps plannings) get(h Hash) *planning {
	if p, ok := ps[h]; ok {
		return &p
	}

	return nil
}

func (ps plannings) add(h Hash, p *planning) {
	ps[h] = *p
}
