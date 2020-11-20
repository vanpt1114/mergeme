package slack

import (
    "fmt"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"

    "github.com/vanpt1114/mergeme/internal/model"
)

func GetMergedBy(projectId int, iid int) (mergedBy model.Block, author model.Block) {
    url := fmt.Sprintf("https://git.teko.vn/api/v4/projects/%d/merge_requests/%d", projectId, iid)
    bearer := "Bearer " + os.Getenv("GITLAB_TOKEN")
    req, err := http.NewRequest(http.MethodGet, url, nil)
    req.Header.Add("Authorization", bearer)
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return
    }
    defer resp.Body.Close()
    var data model.MergeRequestApi
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return
    }

    err = json.Unmarshal([]byte(body), &data)
    if err != nil {
        return
    }

    mergedBy.Type = "section"
    mergedBy.Text = &model.Child{
        Type: "mrkdwn",
        Text: fmt.Sprintf("*Merged by:* <@%s>", data.MergeBy.Username),
    }

    author.Type = "context"
    author.Elements = &[]model.Child{
        {
            Type:       "image",
            ImageUrl:   data.OriAuthor.AvatarUrl,
            AltText:    "author image",
        },
        {
            Type:   "plain_text",
            Text:   data.OriAuthor.Name,
            Emoji:  true,
        },
    }

    return mergedBy, author
}