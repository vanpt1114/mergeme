package message_display

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/vanpt1114/mergeme/internal/bot"
	"github.com/vanpt1114/mergeme/internal/model"

	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

func Open(m model.Message, r, channel string, mr *gitlab.MergeEvent, gl *gitlab.Client, sl *slack.Client, rl *redis.Client) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		m.Description,
		bot.NewGetReviewers(mr.Project.ID, mr.ObjectAttributes.IID, gl),
		m.Footer,
	}

	_, respTS, err := sl.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))

	if err != nil {
		panic(err)
	}

	// Create a redis key with timestamp, so the next event can update to the same thread
	bot.NewUpdateSlackTs(r, respTS, rl)
	bot.NewUpdateSlackTs(fmt.Sprintf("%s:lc", r), mr.ObjectAttributes.LastCommit.ID, rl)
}
