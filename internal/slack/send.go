package slack

import (
    "fmt"
//     "bytes"
//     "encoding/json"
    "context"
//     "net/http"
    "os"

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/vanpt1114/mergeme/config"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var rdb = redis.NewClient(&redis.Options{
    Addr:     os.Getenv("REDIS_HOST"),
    Password: "",
    DB:       0,
})

var url = "https://slack.com/api/chat.postMessage"
var slack_update = "https://slack.com/api/chat.update"
var bearer = "Bearer " + os.Getenv("SLACK_TOKEN")

type Message struct {
    Author      model.Block
    Url         model.Block
    Description model.Block
    Reviewers   model.Block
    Footer      model.Block
}


func SendMessage(m Message, projectId int, objectAttributes model.ObjectAttributes) {
    channel := config.CheckAllow(projectId)
    redisKey := fmt.Sprintf("gitlab:mr:%d", objectAttributes.Id)

    // Switch-case by event `action` field
    switch objectAttributes.Action {
    case "open":
        Open(&m, redisKey, &objectAttributes, channel)
    case "update":
        Update(&m, redisKey, &objectAttributes, channel)
    case "close":
        Close(&m, redisKey, &objectAttributes, channel)
    case "merge":
        Merge(&m, redisKey, &objectAttributes, channel, projectId)
    }
}
