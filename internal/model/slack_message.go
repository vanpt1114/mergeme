package model

import (
    "fmt"
    "strings"
)

func Author (user User) (author Block) {
    author.Type = "context"
    author.Elements = &[]Child{
        {
            Type:       "image",
            ImageUrl:   user.AvatarUrl,
            AltText:    "default alt",
        },
        {
            Type:   "plain_text",
            Text:   user.Name,
            Emoji:  true,
        },
    }
    return author
}

func Url(data ObjectAttributes) (url Block) {
    text := fmt.Sprintf("<%s|*#%d: %s*>\n`%s` ➜ `%s`", data.Url, data.Iid, data.Title, data.SourceBranch, data.TargetBranch)
    url.Type = "section"
    url.Text = &Child{
        Type: "mrkdwn",
        Text: text,
    }
    return url
}

func Description(data ObjectAttributes) (desc Block) {
    desc.Type = "section"
    desc.Text = &Child{
        Type: "mrkdwn",
        Text: "*Description:*\n" + data.Description,
    }
    return desc
}

func Reviewers() (reviewers Block) {
    reviewers.Type = "section"
    reviewers.Text = &Child{
        Type: "mrkdwn",
        Text: "*Reviewers:* " + "<@van.pt> <@van.pt>",
    }
    return reviewers
}

func ReturnAssignees(data []Assignee) Block {
    array := []string{}
    for _, val := range data {
        array = append(array, fmt.Sprintf("<@%s>", val.Username))
    }
    return Block{
        Type: "section",
        Text: &Child{
            Type: "mrkdwn",
            Text: "*Reviewers: *" + strings.Join(array, " "),
        },
    }
}
