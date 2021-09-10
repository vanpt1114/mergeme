package service

import (
    "context"
    "fmt"
    "os"

    "github.com/go-redis/redis/v8"
    "github.com/vanpt1114/mergeme/config"
    "github.com/vanpt1114/mergeme/internal/model"
    //"github.com/vanpt1114/mergeme/internal/service"
    "github.com/xanzy/go-gitlab"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
    Addr:     os.Getenv("REDIS_HOST"),
    Password: "",
    DB:       1,
})

var url = "https://slack.com/api/chat.postMessage"
var slack_update = "https://slack.com/api/chat.update"
var bearer = "Bearer " + os.Getenv("SLACK_TOKEN")
var gitlab_token = os.Getenv("GITLAB_TOKEN")
var gitlab_url = os.Getenv("GITLAB_URL")

const (
	OpenMRColor = "#108548"
    MergedColor = "#1F75CB"
    ClosedColor = "#DD2B0E"
    BotIcon     = ":buff-mr:"
)

type Message struct {
    Author      model.Block
    Url         model.Block
    Description model.Block
    Footer      model.Block
}


func (s *Service) SendMessage(m Message, projectId int, objectAttributes gitlab.MergeEvent) {
    channel := config.CheckAllow(projectId)
    redisKey := fmt.Sprintf("service:mr:%d", objectAttributes.ObjectAttributes.ID)

    // Switch-case by event `action` field
    switch objectAttributes.ObjectAttributes.Action {
    case "open":
        s.Open(&m, redisKey, &objectAttributes, projectId, channel)
    case "update":
        s.Update(&m, redisKey, &objectAttributes, projectId, channel)
    case "close":
        s.Close(&m, redisKey, &objectAttributes, channel)
    case "merge":
        s.Merge(&m, redisKey, &objectAttributes, projectId, channel)
    }
}