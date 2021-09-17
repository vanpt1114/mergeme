package service

import (
    "context"
    "fmt"
    "github.com/slack-go/slack"
    "github.com/vanpt1114/mergeme/config"
    "github.com/xanzy/go-gitlab"
)

var ctx = context.Background()

const (
    BotIcon     = ":buff-mr:"
)

type Message struct {
    Author      slack.Block
    Url         slack.Block
    Description slack.Block
    Footer      slack.Block
}

func (s *Service) SendMessage(m Message, projectId int, mr gitlab.MergeEvent) {
    channel, err := config.CheckAllow(projectId)
    if err != nil {
        panic(err)
    }
    redisKey := fmt.Sprintf("service:mr:%d", mr.ObjectAttributes.ID)

    // Switch-case by event `action` field
    switch mr.ObjectAttributes.Action {
    case "open", "reopen":
        s.Open(m, redisKey, &mr, projectId, channel)
    case "update":
        s.Update(m, redisKey, &mr, projectId, channel)
    case "close":
        s.Close(m, redisKey, &mr, projectId, channel)
    case "merge":
        s.Merge(m, redisKey, &mr, projectId, channel)
    }
}
