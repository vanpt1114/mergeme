package message_display


import (
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/bot"
	"github.com/vanpt1114/mergeme/internal/model"
	"github.com/xanzy/go-gitlab"
)

func Close(m model.Message, r, channel string, mr *gitlab.MergeEvent, gl *gitlab.Client, sl *slack.Client, rl *redis.Client) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		bot.NewGetClosedBy(mr.Project.ID, mr.ObjectAttributes.IID, gl),
	}

	timestamp, err := rl.Get(ctx, r).Result()

	if err == redis.Nil {
		// Redis key does not exists, rarely happen
		// Maybe redis cache have been cleared
		// => Post new message

		// In some case, I think we should skip the "closed" event if the redis key lost,
		// because it not some kind of event that we should acknowledge.

		_, respTS, err := sl.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}

		bot.NewUpdateSlackTs(r, respTS, rl)
	} else if err != nil {
		panic(err)
	} else {
		// Redis key exists, so make an Update to existing thread
		_, _, _, err := sl.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
	}
}
