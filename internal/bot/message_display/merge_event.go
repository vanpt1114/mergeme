package message_display

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/bot"
	"github.com/vanpt1114/mergeme/internal/model"
	"github.com/xanzy/go-gitlab"
)

func Merge(m model.Message, r, channel string, mr *gitlab.MergeEvent, gl *gitlab.Client, sl *slack.Client, rl *redis.Client) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		bot.NewGetMergedBy(mr.Project.ID, mr.ObjectAttributes.IID, gl),
		m.Footer,
	}
	timestamp, err := rl.Get(ctx, r).Result()

	if err == redis.Nil {
		// Redis key does not exist, rarely happen
		fmt.Println("does not process")
	} else if err != nil {
		panic(err)
	} else {
		// Redis key exists, so make an update to the existing thread
		_, _, _, err := sl.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
	}
}
