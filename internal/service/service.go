package service

import (
    "github.com/go-redis/redis/v8"
    "github.com/slack-go/slack"
    "github.com/vanpt1114/mergeme/config"
    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/xanzy/go-gitlab"
    "strconv"
)

type Service struct {
    gitlab  *gitlab.Client
    slack   *slack.Client
    redis   *redis.Client
}

func (s *Service) HandleEvent(data gitlab.MergeEvent) {
    projectID := data.Project.ID
    msgBlock := model.Message{
        Author:      model.Author(data.User),
        Url:         model.Url(data),
        Description: model.Description(data),
        Footer:      model.Footer(data),
    }
    s.SendMessage(msgBlock, projectID, data)
}

func NewService(cfg *config.Config) *Service {
    gl, err := gitlab.NewClient(cfg.Gitlab.Token, gitlab.WithBaseURL(cfg.Gitlab.URL))
    if err != nil {
        panic(err)
    }
    slackClient := slack.New(cfg.SlackToken)
    if err != nil {
        panic(err)
    }
    redisDB, err := strconv.Atoi(cfg.Redis.DB)
    if err != nil {
        panic(err)
    }
    rdb := redis.NewClient(&redis.Options{
        Addr:     cfg.Redis.Host,
        Password: cfg.Redis.Password,
        DB:       redisDB,
    })
    server := &Service{
        gitlab: gl,
        slack: slackClient,
        redis: rdb,
    }
    return server
}

func (s *Service) shouldSkipMergeRequest(event gitlab.MergeEvent) bool {
    // Only allowed selected projects, listed in config/allow.go
	if _, err := config.CheckAllow(event.Project.ID); err != nil {
	    return true
    }

    // Reject WIP/Draft MR
    if event.ObjectAttributes.WorkInProgress {
    	return true
    }

    // Check state of MR, reject merged, closed MR
    if event.ObjectAttributes.State != "opened" && event.ObjectAttributes.Action == "update" {
        return true
    }

    return false
}

//func (s *Service) shouldSkipAction(event slack.InteractionCallback) bool {
//    for _, action := range event.ActionCallback.BlockActions {
//        var tmp model.CustomAction
//        err := json.Unmarshal([]byte(action.Value), &tmp)
//        if err != nil {
//            panic(err)
//        }
//        fmt.Println(tmp)
//
//        if tmp.Clicked {
//            return true
//        }
//    }
//    return false
//}

