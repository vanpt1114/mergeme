package service

import (
    "fmt"
    "github.com/slack-go/slack"
    "github.com/xanzy/go-gitlab"
)

func (s *Service) Open(m Message, r string, o *gitlab.MergeEvent, projectId int, channel string) {
    msgBlock := []slack.Block{
        m.Author,
        m.Url,
        m.Description,
        s.GetReviewers(projectId, o.ObjectAttributes.IID),
        m.Footer,
    }

    slackClient := slack.New(bearer)

    _, respTS, err := slackClient.PostMessage(channel, slack.MsgOptionBlocks(msgBlock...))

    if err != nil {
        panic(err)
    }

    // Create a redis key with timestamp, so the next event can update to the same thread
    UpdateSlackTs(r, respTS)
    UpdateSlackTs(fmt.Sprintf("%s:lc", r), o.ObjectAttributes.LastCommit.ID)
}
