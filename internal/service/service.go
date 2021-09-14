package service

import (
    "fmt"
    "github.com/slack-go/slack"
    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/xanzy/go-gitlab"
    "strings"
)

type Service struct {
    gitlab  *gitlab.Client
    slack   *slack.Client
}

func (s *Service) Handle(data gitlab.MergeEvent) {
    projectID := data.Project.ID
    msgBlock := Message{
        Author:      model.Author(data.User),
        Url:         model.Url(data),
        Description: model.Description(data),
        Footer:      model.Footer(data),
    }
    s.SendMessage(msgBlock, projectID, data)
}

func NewService(token, baseURL, slackToken string) *Service {
    gl, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
    if err != nil {
        panic(err)
    }
    slackClient := slack.New(slackToken)
    if err != nil {
        panic(err)
    }
    server := &Service{
        gitlab: gl,
        slack: slackClient,
    }
    return server
}

func UpdateSlackTs(r, ts string) {
    err := rdb.Set(ctx, r, ts, 0).Err()
    if err != nil {
        panic(err)
    }
}

func (s *Service) GetReviewers(projectID, mrIID int) (reviewers *slack.SectionBlock) {
    mr, _, err := s.gitlab.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
    if err != nil {
        panic(err)
    }

    array := []string{}
    for _, val := range mr.Reviewers {
        array = append(array, fmt.Sprintf("<@%s>", val.Username))
    }

    reviewersText := slack.NewTextBlockObject("mrkdwn",
        fmt.Sprintf("*Reviewers: *" + strings.Join(array, " ")), false, false)

    reviewers = slack.NewSectionBlock(reviewersText, nil, nil)

    return reviewers
}

func (s *Service) GetMergedBy(projectID, mrIID int) (mergedBy *slack.SectionBlock) {
    mr, _, err := s.gitlab.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
    if err != nil {
        panic(err)
    }
    textBlock := slack.NewTextBlockObject(
        "mrkdwn",
        fmt.Sprintf("*Merged by: *<@%s>", mr.MergedBy.Username),
        false,
        false)
    mergedBy = slack.NewSectionBlock(textBlock, nil, nil)

    return mergedBy
}

func (s *Service) GetClosedBy(projectID, mrIID int) (closedBy *slack.SectionBlock) {
    mr, _, err := s.gitlab.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
    if err != nil {
        panic(err)
    }
    textBlock := slack.NewTextBlockObject(
        "mrkdwn",
        fmt.Sprintf("*Closed by: *<@%s>", mr.ClosedBy.Username),
        false,
        false)
    closedBy = slack.NewSectionBlock(textBlock, nil, nil)

    return closedBy
}
