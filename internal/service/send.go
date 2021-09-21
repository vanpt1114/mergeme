package service

import (
    "context"
    "fmt"
    "github.com/slack-go/slack"
    "github.com/vanpt1114/mergeme/config"
    "github.com/vanpt1114/mergeme/internal/bot/message_display"
    "github.com/vanpt1114/mergeme/internal/model"
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

func (s *Service) SendMessage(m model.Message, projectId int, mr gitlab.MergeEvent) {
    channel, err := config.CheckAllow(projectId)
    if err != nil {
        panic(err)
    }
    redisKey := fmt.Sprintf("service:mr:%d", mr.ObjectAttributes.ID)

    // Switch-case by event `action` field
    switch mr.ObjectAttributes.Action {
    case "open", "reopen":
        //s.Open(m, redisKey, &mr, projectId, channel)
        message_display.Open(m, redisKey, channel, &mr, s.gitlab, s.slack, s.redis)
    case "update":
        message_display.Update(m, redisKey, channel, &mr, s.gitlab, s.slack, s.redis)
    case "close":
        message_display.Close(m, redisKey, channel, &mr, s.gitlab, s.slack, s.redis)
    case "merge":
        message_display.Merge(m, redisKey, channel, &mr, s.gitlab, s.slack, s.redis)
    default:
        fmt.Printf("%s does not supported yet", mr.ObjectAttributes.Action)
    }
}
