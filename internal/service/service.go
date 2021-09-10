package service

import (
    "fmt"
    "github.com/vanpt1114/mergeme/internal/model"
    "strings"

    "github.com/xanzy/go-gitlab"
)

type Service struct {
    gitlab  *gitlab.Client
}

func (s *Service) Handle(data gitlab.MergeEvent) {
    projectId := data.Project.ID
    result := Message{
        Author:         model.Author(data.User),
        Url:            model.Url(data),
        Description:    model.Description(data),
        Footer:         model.Repo(data),
    }
    s.SendMessage(result, projectId, data)
}

func NewService(token, baseURL string) *Service {
    gl, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL))
    if err != nil {
        panic(err)
    }
    server := &Service{
        gitlab: gl,
    }
    return server
}

func (s *Service) GetReviewers(projectID, mrIID int) (reviewers model.Block) {
   mr, _, err := s.gitlab.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
   if err != nil {
       panic(err)
   }

   array := []string{}
   for _, val := range mr.Reviewers {
       array = append(array, fmt.Sprintf("<@%s>", val.Username))
   }

   reviewers.Type = "section"
   reviewers.Text = &model.Child{
       Type: "mrkdwn",
       Text: "*Reviewers: *" + strings.Join(array, " "),
   }
   return reviewers
}

func (s *Service) GetMergedBy(projectID, mrIID int) (mergedBy, author model.Block) {
   mr, _, err := s.gitlab.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
   if err != nil {
       panic(err)
   }

   mergedBy.Type = "section"
   mergedBy.Text = &model.Child{
       Type: "mrkdwn",
       Text: fmt.Sprintf("*Merged by:* <@%s>", mr.MergedBy.Username),
   }

    author.Type = "context"
    author.Elements = &[]model.Child{
        {
            Type:       "image",
            ImageUrl:   mr.Author.AvatarURL,
            AltText:    "author image",
        },
        {
            Type:   "plain_text",
            Text:   mr.Author.Name,
            Emoji:  true,
        },
    }
   return mergedBy, author
}