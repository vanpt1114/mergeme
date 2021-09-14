package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

func (s *Service) Merge(m Message, r string, mr *gitlab.MergeEvent, projectID int, channel string) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		s.GetMergedBy(projectID, mr.ObjectAttributes.IID),
		m.Footer,
	}

	slackClient := slack.New(bearer)
	timestamp, err := rdb.Get(ctx, r).Result()

	if err == redis.Nil {
		// Redis key does not exist, rarely happen
		_, respTS, err := slackClient.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
		UpdateSlackTs(r, respTS)
	} else if err != nil {
		panic(err)
	} else {
		// Redis key exists, so make an update to the existing thread
		_, _, _, err := slackClient.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
	}
}
