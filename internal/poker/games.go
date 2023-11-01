package poker

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	gameHash string

	games map[gameHash]game
	game  struct {
		votes     votes
		createdAt int
		updatedAt int
		creator   *tgbotapi.User
		jiraIssue string
	}

	voter *tgbotapi.User
)

func newHashFromBotCommandMessage(m *tgbotapi.Message, jiraIssue string) gameHash {
	if m == nil {
		return ""
	}
	md5Hasher := md5.New()
	md5Hasher.Write([]byte(fmt.Sprintf(
		"%d:%d:%s",
		m.From.ID,
		m.Chat.ID,
		jiraIssue,
	)))

	return gameHash(hex.EncodeToString(md5Hasher.Sum(nil)))
}

func newGame(j string, m *tgbotapi.Message) *game {
	return &game{createdAt: m.Date, updatedAt: m.Date, creator: m.From, jiraIssue: j, votes: make(votes, 0)}
}

func (ps games) get(h gameHash) *game {
	if p, ok := ps[h]; ok {
		return &p
	}

	return nil
}

func (ps games) add(h gameHash, p *game) {
	ps[h] = *p
}

func (p game) addVote(v vote) {
	p.votes.add(v)
}

func (p game) listVoters() []voter {
	voters := make([]voter, 0)
	for _, v := range p.votes {
		voters = append(voters, v.voter)
	}

	return voters
}
