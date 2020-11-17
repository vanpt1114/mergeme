package gitlab

import (

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/vanpt1114/mergeme/internal/slack"
)

func Handle(data model.MergeRequest) {
    projectId := data.Project.Id
    objectAttributes := data.ObjectAttributes
    result := slack.Message{
        Author:         model.Author(data.User),
        Url:            model.Url(data.ObjectAttributes),
        Description:    model.Description(data.ObjectAttributes),
        Reviewers:      model.ReturnAssignees(data.Assignees),
    }
    slack.SendMessage(result, projectId, objectAttributes)
}