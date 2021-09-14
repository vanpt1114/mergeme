package service

import (
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

func (s *Service) Close(m Message, r string, o *gitlab.MergeEvent, projectID int, channel string) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		s.GetClosedBy(projectID, o.ObjectAttributes.IID),
	}

	timestamp, err := rdb.Get(ctx, r).Result()

	if err == redis.Nil {
		// Redis key does not exists, rarely happen
		// Maybe redis cache have been cleared
		// => Post new message
		_, respTS, err := s.slack.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}

		UpdateSlackTs(r, respTS)
	} else if err != nil {
		panic(err)
	} else {
		// Redis key exists, so make an Update to existing thread
		_, _, _, err := s.slack.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
	}
}
