package service

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
)

func (s *Service) Update(m Message, r string, mr *gitlab.MergeEvent, projectID int, channel string) {
	msgBlock := []slack.Block{
		m.Author,
		m.Url,
		m.Description,
		s.GetReviewers(projectID, mr.ObjectAttributes.IID),
		m.Footer,
	}

	timestamp, err := rdb.Get(ctx, r).Result()

	if err == redis.Nil {
		// "[Update_1] Redis key doesn't exist"
		_, respTS, err := s.slack.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))

		if err != nil {
			panic(err)
		}
		UpdateSlackTs(r, respTS)
		UpdateSlackTs(fmt.Sprintf("%s:lc", r), mr.ObjectAttributes.LastCommit.ID)
	} else if err != nil {
		panic("[Update_1] Err")
	} else {
		// [Update_1] Found redis key
		lastCommit, err := rdb.Get(ctx, fmt.Sprintf("%s:lc", r)).Result()
		if err == redis.Nil {
			fmt.Println("[Update_2] last_commit not found, cache may be cleared")
		} else if err != nil {
			panic("[Update_2] Err")
		} else {
			if mr.ObjectAttributes.LastCommit.ID == lastCommit {
				// last_commit does not change, so re-update thread"
				_, _, _, err := s.slack.UpdateMessage(channel, timestamp, slack.MsgOptionBlocks(msgBlock...))
				if err != nil {
					panic(err)
				}
			} else {
				// last_commit is different, so make a post with sub-message"
				_, _, err := s.slack.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...), slack.MsgOptionTS(timestamp))
				if err != nil {
					panic(err)
				}
			}
		}
	}
}
