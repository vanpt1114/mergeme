package model

import (
    "fmt"
    "github.com/slack-go/slack"
    "github.com/xanzy/go-gitlab"
    "html"
    "regexp"
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

func Author(user *gitlab.EventUser) (author *slack.ContextBlock) {
	textBlock := slack.NewTextBlockObject(
	    "plain_text",
	    fmt.Sprintf("<@%s>", user.Username),false,false)
	author = slack.NewContextBlock("", []slack.MixedElement{textBlock}...)
	return author
}

func Url(data gitlab.MergeEvent) (url *slack.SectionBlock) {
    textBlock :=  slack.NewTextBlockObject("mrkdwn",
        fmt.Sprintf(
            "<%s|*#%d: %s*>\n`%s` ➜ `%s`",
            data.ObjectAttributes.URL,
            data.ObjectAttributes.IID,
            html.UnescapeString(data.ObjectAttributes.Title),
            data.ObjectAttributes.SourceBranch,
            data.ObjectAttributes.TargetBranch,
        ), false, false)
    url = slack.NewSectionBlock(textBlock, nil, nil)
    return url
}

func Description(data gitlab.MergeEvent) (desc *slack.SectionBlock) {
    description := SlackMarkDown(data.ObjectAttributes.Description)
	textBlock := slack.NewTextBlockObject("mrkdwn", description + " ", false, false)
	desc = slack.NewSectionBlock(textBlock, nil, nil)
    return desc
}

func Footer(mr gitlab.MergeEvent) (footer *slack.ContextBlock) {
    textBlock := slack.NewTextBlockObject("plain_text", mr.Project.Name, false, false)
    if len(mr.Project.AvatarURL) != 0 {
        repoIconBlock := slack.NewImageBlockElement(mr.Project.AvatarURL, "Repository icon")
        footer = slack.NewContextBlock("", []slack.MixedElement{repoIconBlock, textBlock}...)
        return footer
    }
    return slack.NewContextBlock("", []slack.MixedElement{textBlock}...)
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
