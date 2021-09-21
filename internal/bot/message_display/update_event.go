package message_display


import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/bot"
	"github.com/vanpt1114/mergeme/internal/model"
	"github.com/xanzy/go-gitlab"
)

var ctx = context.Background()

func Update(m model.Message, r, channel string, mr *gitlab.MergeEvent, gl *gitlab.Client, sl *slack.Client, rl *redis.Client) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		m.Description,
		bot.NewGetReviewers(mr.Project.ID, mr.ObjectAttributes.IID, gl),
		m.Footer,
	}

	timestamp, err := rl.Get(ctx, r).Result()

	if err == redis.Nil {
		// "[Update_1] Redis key doesn't exist"
		_, respTS, err := sl.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))

		if err != nil {
			panic(err)
		}
		bot.NewUpdateSlackTs(r, respTS, rl)
		bot.NewUpdateSlackTs(fmt.Sprintf("%s:lc", r), mr.ObjectAttributes.LastCommit.ID, rl)
	} else if err != nil {
		panic(err)
	} else {
		// [Update_1] Found redis key
		lastCommit, err := rl.Get(ctx, fmt.Sprintf("%s:lc", r)).Result()
		if err == redis.Nil {
			fmt.Println("[Update_2] last_commit not found, cache may be cleared")
		} else if err != nil {
			panic("[Update_2] Err")
		} else {
			if mr.ObjectAttributes.LastCommit.ID == lastCommit {
				// last_commit does not change, so re-update thread"
				_, _, _, err := sl.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
				if err != nil {
					panic(err)
				}
			} else {
				// last_commit is different, so make a post with sub-message"
				_, _, err := sl.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...), slack.MsgOptionTS(timestamp))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
