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
	timestamp, err := s.redis.Get(ctx, r).Result()

	if err == redis.Nil {
		// Redis key does not exist, rarely happen
		_, respTS, err := s.slack.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
		s.UpdateSlackTs(r, respTS)
	} else if err != nil {
		panic(err)
	} else {
		// Redis key exists, so make an update to the existing thread
		_, _, _, err := s.slack.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
		if err != nil {
			panic(err)
		}
	}
}
