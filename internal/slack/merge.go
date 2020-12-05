package slack

import (
//     "fmt"
    "net/http"
    "bytes"
    "encoding/json"

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/go-redis/redis/v8"
)

func Merge(m *Message, r string, o *model.ObjectAttributes, channel string, projectId int) {
    mergedBy, author := GetMergedBy(projectId, o.Iid)
    dataAttachments := []model.Attachment{
        model.Attachment{
            Color: "#030321",
            Blocks: []model.Block{
                author,
                m.Url,
//                 m.Description,
                mergedBy,
                m.Footer,
            },
        },
    }

    timestamp, err := rdb.Get(ctx, r).Result()
    if err == redis.Nil {
        // Redis key does not exist, rarely happen
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
        UpdateSlackTs(r, ts)
    } else if err != nil {
        panic(err)
    } else {
        // Redis key exists, so make an update to the existing thread
        dataToSend, _ := json.Marshal(&model.SlackPayload{
            Channel:        channel,
            Ts:             timestamp,
            Username:       "MergeMe",
            IconEmoji:      ":buff-mr:",
            Attachments:    dataAttachments,
        })
        req, err := http.NewRequest(http.MethodPost, slack_update, bytes.NewBuffer(dataToSend))
        req.Header.Add("Authorization", bearer)
        req.Header.Add("Content-Type", "application/json")
        client := &http.Client{}

        resp, err := client.Do(req)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()
    }
}
