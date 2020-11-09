package slack

import (
    "fmt"
//     "bytes"
//     "net/http"
//     "encoding/json"

    model   "../model"
)

type Child struct {
    Type     string  `json:"type,omitempty"`
    ImageUrl string  `json:"image_url,omitempty"`
    AltText  string  `json:"alt_text,omitempty"`
    Text     string  `json:"text,omitempty"`
    Emoji    bool    `json:"emoji,omitempty"`
}


type Block struct {
    Type        string      `json:"type"`
    Elements    *[]Child    `json:"elements,omitempty"`
    Text        *Child      `json:"text,omitempty"`
}

type Attachment struct {
    Color   string  `json:"color"`
    Blocks  []Block `json:"blocks"`
}

type SlackPayload struct {
    Channel         string          `json:"channel"`
    Username        string          `json:"username"`
    IconEmoji       string          `json:"icon_emoji"`
    Attachments     []Attachment    `json:"attachments"`
}

func Author (user model.User) (author Block) {
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

func Url(data model.ObjectAttributes) (url Block) {
    text := fmt.Sprintf("<%s|*#%d: %s*>\n`%s` ➜ `%s`", data.Url, data.Iid, data.Title, data.SourceBranch, data.TargetBranch)
    url.Type = "section"
    url.Text = &Child{
        Type: "mrkdwn",
        Text: text,
    }
    return url
}

func Description(data model.ObjectAttributes) (desc Block) {
    desc.Type = "section"
    desc.Text = &Child{
        Type: "mrkdwn",
        Text: "*Description:*\n" + data.Description,
    }
    return desc
}

