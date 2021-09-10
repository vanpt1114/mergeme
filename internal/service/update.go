package service

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/xanzy/go-gitlab"
    "net/http"

    "github.com/go-redis/redis/v8"
    "github.com/vanpt1114/mergeme/internal/model"
)

func (s *Service) Update(m *Message, r string, o *gitlab.MergeEvent, projectId int, channel string) {
    reviewers := s.GetReviewers(projectId, o.ObjectAttributes.IID)
    fmt.Println(reviewers)
    timestamp, err := rdb.Get(ctx, r).Result()
    if err == redis.Nil {
        // "[Update_1] Redis key doesn't exist"
        dataAttachments := []model.Attachment{
            {
                Color: OpenMRColor,
                Blocks: []model.Block{
                    m.Author,
                    m.Url,
                    m.Description,
                    reviewers,
                    m.Footer,
                },
            },
        }

        dataToSend, _ := json.Marshal(&model.SlackPayload{
            Channel:     channel,
            Username:    "MergeMe",
            IconEmoji:   BotIcon,
            Attachments: dataAttachments,
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
        UpdateSlackTs(fmt.Sprintf("%s:lc", r), o.ObjectAttributes.LastCommit.ID)
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
            if o.ObjectAttributes.LastCommit.ID == lastCommit {
                fmt.Println("[Update_2] last_commit does not change, so re-update thread")
                dataAttachments := []model.Attachment{
                    {
                        Color: OpenMRColor,
                        Blocks: []model.Block{
                            m.Author,
                            m.Url,
                            m.Description,
                            reviewers,
                            m.Footer,
                        },
                    },
                }

                dataToSend, _ := json.Marshal(&model.SlackPayload{
                    Channel:     channel,
                    Ts:          timestamp,
                    Username:    "MergeMe",
                    IconEmoji:   BotIcon,
                    Attachments: dataAttachments,
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
                    {
                        Color: OpenMRColor,
                        Blocks: []model.Block{
                            m.Author,
                            reviewers,
                        },
                    },
                }

                dataToSend, _ := json.Marshal(&model.SlackPayload{
                    Channel:     channel,
                    ThreadTs:    timestamp,
                    Username:    "MergeMe",
                    IconEmoji:   BotIcon,
                    Attachments: dataAttachments,
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
                //UpdateSlackTs(fmt.Sprintf("%s:lc", r), o.LastCommit.Id)
            }
        }
    }
}
