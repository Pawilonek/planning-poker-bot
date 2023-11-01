package poker

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type (
	votes map[int64]vote
	vote  struct {
		voter     voter
		value     string
		createdAt int
		updatedAt int
	}
)

func newVote(u *tgbotapi.User, d callbackData) vote {
	return vote{
		voter:     u,
		value:     d.data,
		createdAt: 0,
		updatedAt: 0,
	}
}

func (vs votes) add(v vote) {
	vs[v.voter.ID] = v
}
