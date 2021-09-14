package service

import (
    "context"
    "fmt"
    "github.com/slack-go/slack"
    "os"

    "github.com/go-redis/redis/v8"
    "github.com/vanpt1114/mergeme/config"
    "github.com/xanzy/go-gitlab"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
    Addr:     os.Getenv("REDIS_HOST"),
    Password: "",
    DB:       1,
})

var bearer = os.Getenv("SLACK_TOKEN")
var gitlab_token = os.Getenv("GITLAB_TOKEN")
var gitlab_url = os.Getenv("GITLAB_URL")

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
    channel := config.CheckAllow(projectId)
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
