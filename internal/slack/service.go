package slack

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "strings"

    "github.com/vanpt1114/mergeme/internal/model"
)

var GITLAB_URL = os.Getenv("GITLAB_URL")

func GetMergedBy(projectId int, iid int) (mergedBy model.Block, author model.Block) {
    url := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d", GITLAB_URL, projectId, iid)
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
            ImageUrl:   data.OriAuthor.AvatarURL,
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

func UpdateSlackTs(r, ts string) {
    err := rdb.Set(ctx, r, ts, 0).Err()
    if err != nil {
        panic(err)
    }
}

func DecodeSlackResponse(r *http.Response) string {
    decoder := json.NewDecoder(r.Body)
    var t model.SlackResponsePayload
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    if t.Ok != true {
        panic(t.Error)
    }
    return t.Ts
}

func GetReviewers(projectId int, iid int) (reviewers model.Block) {
    fmt.Println("Execute GetReviewers function ...")
    url := fmt.Sprintf("%s/api/v4/projects/%d/merge_requests/%d", GITLAB_URL, projectId, iid)
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
    array := []string{}
    for _, val := range data.Reviewers {
        fmt.Println(val)
        array = append(array, fmt.Sprintf("<@%s>", val.Username))
    }

    reviewers.Type = "section"
    reviewers.Text = &model.Child{
        Type: "mrkdwn",
        Text: "*Reviewers: *" + strings.Join(array, " "),
    }

    return reviewers
}
