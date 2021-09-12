package model

import (
    "fmt"
	"github.com/xanzy/go-gitlab"
	"html"
    "regexp"

    "github.com/slack-go/slack"
)

const (
    toDoIcon        = ":todo:"
    toDoDoneIcon    = ":todo_done:"
)

var Pattern = map[string]string{
    `\-\s\[x\]`:    toDoDoneIcon,
    `\-\s\[\s\]`:   toDoIcon,
    `\-\s\[\+\]`:   toDoDoneIcon,
    `\*\s\[x\]`:    toDoDoneIcon,
    `\*\s\[\s\]`:   toDoIcon,
}

var Bold = map[string]string{
    `\*\*\*([0-9a-zA-Z-./#: ]+)\*\*\*`: "*_", // Bold & Italic
    `\*\*([0-9a-zA-Z-./#: ]+)\*\*`: "*", // Bold
    `\_\_([0-9a-zA-Z-./#: ]+)\_\_`: "*", // Bold
    `#{1,6}\s([0-9a-zA-Z-./#: ]+)`: "*", // Heading
}

func Reverse(s string) string {
    runes := []rune(s)
    for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
        runes[i], runes[j] = runes[j], runes[i]
    }
    return string(runes)
}

func NewAuthor(user *gitlab.EventUser) (author slack.Blocks) {
    author.BlockSet = []slack.Block{
        slack.ImageBlock{
            Type:       "image",
            ImageURL:   user.AvatarURL,
            AltText:    "default alt",
        },
        slack.TextBlockObject{
            Type:       "plain_text",
            Text:       fmt.Sprintf("<@%s>", user.Username),
            Emoji:      true,
        },
    }
    return author
}

func Author(user *gitlab.EventUser) (author Block) {
    author.Type = "context"
    author.Elements = &[]Child{
        {
            Type:       "image",
            ImageUrl:   user.AvatarURL,
            AltText:    "default alt",
        },
        {
            Type:   "plain_text",
            Text:   fmt.Sprintf("<@%s>", user.Username),
            Emoji:  true,
        },
    }
    return author
}

func Url(data gitlab.MergeEvent) (url Block) {
    text := fmt.Sprintf(
        "<%s|*#%d: %s*>\n`%s` ➜ `%s`",
        data.ObjectAttributes.URL,
        data.ObjectAttributes.IID,
        html.UnescapeString(data.ObjectAttributes.Title),
        data.ObjectAttributes.SourceBranch,
        data.ObjectAttributes.TargetBranch,
        )
    
    url.Type = "section"
    url.Text = &Child{
        Type: "mrkdwn",
        Text: text,
    }
    return url
}

func Description(data gitlab.MergeEvent) (desc Block) {
    description := SlackMarkDown(data.ObjectAttributes.Description)
    desc.Type = "section"
    desc.Text = &Child{
        Type: "mrkdwn",
        Text: description + " ",
    }
    return desc
}

func Repo(data gitlab.MergeEvent) Block {
    if data.Project.AvatarURL== "" {
        return Block{
            Type: "context",
            Elements: &[]Child{
                {
                    Type:   "plain_text",
                    Text:   data.Project.Name,
                    Emoji:  true,
                },
            },
        }
    }
    return Block{
        Type: "context",
        Elements: &[]Child{
            {
                Type:       "image",
                ImageUrl:   data.Project.AvatarURL,
                AltText:    "default alt",
            },
            {
                Type:   "plain_text",
                Text:   data.Project.Name,
                Emoji:  true,
            },
        },
    }
}

func SlackMarkDown(data string) string {
    for k, v := range Pattern {
        re := regexp.MustCompile(k)
        data = re.ReplaceAllString(data, v)
    }

    // Replace markdown Bold, Italic
    for k, v := range Bold {
        re := regexp.MustCompile(k)
        tmp := re.FindAllStringSubmatch(data, -1)
        for _, val := range tmp {
            re := regexp.MustCompile(regexp.QuoteMeta(val[0]))
            data = re.ReplaceAllString(data, v + val[1] + Reverse(v))
        }
    }

    // Replace markdown Hyperlink
    re := regexp.MustCompile(`\[([a-zA-Z0-9!@#$%^&*.\-/?=+]+)\]\((https?:\/\/[a-zA-Z0-9!@#$%^&*.\-/?=+]+)\)`)
	result := re.FindAllStringSubmatch(data, -1)
	for _, val := range result {
		re := regexp.MustCompile(regexp.QuoteMeta(val[0]))
		hyperlink := fmt.Sprintf("<%s|%s>", val[2], val[1])
		data = re.ReplaceAllString(data, hyperlink)
	}
    return data
}
