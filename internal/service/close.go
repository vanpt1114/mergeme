package service

import (
    "bytes"
    "encoding/json"
    "github.com/xanzy/go-gitlab"
    "net/http"

    "github.com/go-redis/redis/v8"
    "github.com/vanpt1114/mergeme/internal/model"
)



func (s *Service) Close(m *Message, r string, o *gitlab.MergeEvent, channel string) {
    var closeBlock model.Block
    closeBlock.Type = "section"
    closeBlock.Text = &model.Child{
        Type: "mrkdwn",
        Text: "*Closed*",
    }
    dataAttachments := []model.Attachment{
        {
            Color: ClosedColor,
            Blocks: []model.Block{
                m.Author,
                m.Url,
                m.Description,
                closeBlock,
                m.Footer,
            },
        },
    }

    timestamp, err := rdb.Get(ctx, r).Result()
    if err == redis.Nil {
        // Redis key does not exists => Post new message
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
