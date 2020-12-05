package slack

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/go-redis/redis/v8"
)



func Close(m *Message, r string, o *model.ObjectAttributes, channel string) {
    fmt.Println("Did i reach here?")
    var closeBlock model.Block
    closeBlock.Type = "section"
    closeBlock.Text = &model.Child{
        Type: "mrkdwn",
        Text: "*Closed*",
    }
    dataAttachments := []model.Attachment{
        model.Attachment{
            Color: "#1542e6",
            Blocks: []model.Block{
                m.Author,
                m.Url,
                m.Description,
                closeBlock,
//                 m.Footer,
            },
        },
    }

    timestamp, err := rdb.Get(ctx, r).Result()
    if err == redis.Nil {
        // Redis key does not exists => Post new message
        fmt.Println("redis key does not exists")
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

        ts := DecodeSlackResponse(resp)
        UpdateSlackTs(r, ts)
    } else if err != nil {
        panic(err)
    } else {
        // Redis key exists, so make an Update to existing thread
        fmt.Println("redis key does exists")
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
