package model

import (
    "fmt"
    "strings"
    "regexp"
)

var Pattern = map[string]string{
    `\-\s\[x\]`:    ":todo_done:",
    `\-\s\[\s\]`:   ":todo:",
    `\-\s\[\+\]`:   ":todo_done:",
}

var Bold = map[string]string{
    `\*\*\*([0-9a-zA-Z-./#: ]+)\*\*\*`: "*_", // Bold & Italic
//     `\*([0-9a-zA-Z-./#: ]+)\*`: "_", // Italic
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

func Author(user User) (author Block) {
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
    description := SlackMarkDown(data.Description)
    desc.Type = "section"
    desc.Text = &Child{
        Type: "mrkdwn",
//         Text: "*Description:*\n" + description,
        Text: description,
    }
    return desc
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


func Repo(data Project) Block {
    if data.AvatarUrl == "" {
        return Block{
            Type: "context",
            Elements: &[]Child{
                {
                    Type:   "plain_text",
                    Text:   data.Name,
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
                ImageUrl:   data.AvatarUrl,
                AltText:    "default alt",
            },
            {
                Type:   "plain_text",
                Text:   data.Name,
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
    re := regexp.MustCompile(`\[([\s0-9a-zA-Z-.]+)\]\((https?:\/\/[0-9a-zA-Z\-./#]+)\)`)
	result := re.FindAllStringSubmatch(data, -1)
	for _, val := range result {
		re := regexp.MustCompile(regexp.QuoteMeta(val[0]))
		hyperlink := fmt.Sprintf("<%s|%s>", val[2], val[1])
		data = re.ReplaceAllString(data, hyperlink)
	}
    return data
}
