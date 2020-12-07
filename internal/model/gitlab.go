package model

import (
    "fmt"
)

type Commit struct {
    ShortId     string  `json:"short_id"`
    Message     string  `json:"message"`
}

type User struct {
    Name        string  `json:"name"`
    User        string  `json:"username"`
    AvatarUrl   string  `json:"avatar_url"`
    Email       string  `json:"email"`
}

type Project struct {
    Id          int     `json:"id"`
    AvatarUrl   string  `json:"avatar_url"`
    Name        string  `json:"name"`
    Path        string  `json:"path_with_namespace"`
}

type LastCommit struct {
    Id      string  `json:"id"`
    Message string  `json:message`
}

type Assignee struct {
    Username    string  `json:"username"`
}

type ObjectAttributes struct {
    AssigneeId      int         `json:"assignee_id"`
    AuthorId        int         `json:"author_id"`
    CreatedAt       string      `json:"created_at"`
    Description     string      `json:"description"`
    Id              int         `json:"id"`
    Iid             int         `json:"iid"`
    SourceBranch    string      `json:"source_branch"`
    TargetBranch    string      `json:"target_branch"`
    Title           string      `json:"title"`
    Url             string      `json:"url"`
    LastCommit      LastCommit  `json:"last_commit"`
    State           string      `json:"state"`
    Action          string      `json:"action"`
    WIP             bool        `json:"work_in_progress"`
}

type MergeRequest struct {
    ObjectKind          string              `json:"object_kind"`
    EventType           string              `json:"event_type"`
    User                User                `json:"user"`
    Project             Project             `json:"project"`
    ObjectAttributes    ObjectAttributes    `json:"object_attributes"`
    Assignees           []Assignee          `json:"assignees"`
}


func Handle(data MergeRequest) {
    fmt.Println(data)
}
