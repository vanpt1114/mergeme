package slack

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"

    "github.com/vanpt1114/mergeme/internal/model"
)

func Open(m *Message, r string, o *model.ObjectAttributes, channel string) {
    dataAttachments := []model.Attachment{
        model.Attachment{
            Color: OpenMRColor,
            Blocks: []model.Block{
                m.Author,
                m.Url,
                m.Description,
                m.Reviewers,
                m.Footer,
            },
        },
    }

    dataToSend, _ := json.Marshal(&model.SlackPayload{
        Channel:        channel,
        Username:       "MergeMe",
        IconEmoji:      ":buff-mr:",
        Attachments:    dataAttachments,
    })

    req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(dataToSend))

    req.Header.Add("Authorization", bearer)
    req.Header.Add("Content-Type", "application/json")
    client := &http.Client{}

    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    ts := DecodeSlackResponse(resp)

    // Create a redis key with timestamp, so the next event can update to the same thread
    UpdateSlackTs(r, ts)
    UpdateSlackTs(fmt.Sprintf("%s:lc", r), o.LastCommit.Id)
}
