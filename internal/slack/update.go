package slack

import (
    "fmt"
    "net/http"
    "bytes"
    "encoding/json"

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/go-redis/redis/v8"
)

func Update(m *Message, r string, o *model.ObjectAttributes, channel string) {
    timestamp, err := rdb.Get(ctx, r).Result()
    if err == redis.Nil {
        // "[Update_1] Redis key doesn't exist"
        dataAttachments := []model.Attachment{
            model.Attachment{
                Color: "#1542e6",
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
            Channel:        "C019QRNJ6DN",
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
        UpdateSlackTs(fmt.Sprintf("%s:lc", r), o.LastCommit.Id)
    } else if err != nil {
        panic("[Update_1] Err")
    } else {
        // [Update_1] Found redis key
        lastCommit, err := rdb.Get(ctx, fmt.Sprintf("%s:lc", r)).Result()
        if err == redis.Nil {
            fmt.Println("[Update_2] last_commit not found")
        } else if err != nil {
            panic("[Update_2] Err")
        } else {
            if o.LastCommit.Id == lastCommit {
                fmt.Println("[Update_2] last_commit does not change, so re-update thread")
                dataAttachments := []model.Attachment{
                    model.Attachment{
                        Color: "#1542e6",
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
                DecodeSlackResponse(resp)
            } else {
                fmt.Println("[Update_2] last_commit is different, so make a post with sub-message")
                dataAttachments := []model.Attachment{
                    model.Attachment{
                        Color: "#1542e6",
                        Blocks: []model.Block{
                            m.Author,
                            m.Reviewers,
                        },
                    },
                }

                dataToSend, _ := json.Marshal(&model.SlackPayload{
                    Channel:        channel,
                    ThreadTs:       timestamp,
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
                DecodeSlackResponse(resp)
            }
        }
    }
}