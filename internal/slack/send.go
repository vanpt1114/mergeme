package slack

import (
    "fmt"
    "bytes"
    "encoding/json"
    "context"
    "net/http"
    "os"

    "github.com/vanpt1114/mergeme/internal/model"
    "github.com/vanpt1114/mergeme/config"
    "github.com/go-redis/redis/v8"
)

var ctx = context.Background()

type Message struct {
    Author      model.Block
    Url         model.Block
    Description model.Block
    Reviewers   model.Block
    Footer      model.Block
}

func SendMessage(m Message, projectId int, objectAttributes model.ObjectAttributes) {
    channel := config.CheckAllow(projectId)
    rdb := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST"),
        Password: "",
        DB:       0,
    })

    url := "https://slack.com/api/chat.postMessage"
    bearer := "Bearer " + os.Getenv("SLACK_TOKEN")

    redisKey := fmt.Sprintf("gitlab:mr:%d", objectAttributes.Id)

    switch objectAttributes.Action {
    case "open":
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

        decoder := json.NewDecoder(resp.Body)
        var t model.SlackResponsePayload
        err = decoder.Decode(&t)
        if err != nil {
            panic(err)
        }
        defer resp.Body.Close()

        err = rdb.Set(ctx, redisKey, t.Ts, 0).Err()
        if err != nil {
            panic(err)
        }

        lastCommit := objectAttributes.LastCommit.Id
        err = rdb.Set(ctx, fmt.Sprintf("%s:lc", redisKey), lastCommit, 0).Err()
        if err != nil {
            panic(err)
        }
    case "update":
//         fmt.Println("update")
        ts, err := rdb.Get(ctx, redisKey).Result()
        if err == redis.Nil {
            // Update event with no redis key
            fmt.Println("Update event with no redis key")
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

            decoder := json.NewDecoder(resp.Body)
            var t model.SlackResponsePayload
            err = decoder.Decode(&t)
            if err != nil {
                panic(err)
            }
            defer resp.Body.Close()

            err = rdb.Set(ctx, redisKey, t.Ts, 0).Err()
            if err != nil {
                panic(err)
            }

            lastCommit := objectAttributes.LastCommit.Id
            err = rdb.Set(ctx, fmt.Sprintf("%s:lc", redisKey), lastCommit, 0).Err()
            if err != nil {
                panic(err)
            }
        } else {
            // Update event with new `last_commit` id (redis key exists)
            fmt.Println("Update event with new `last_commit` id (redis key exists)")
            lastCommit := objectAttributes.LastCommit.Id
            lc, err := rdb.Get(ctx, fmt.Sprintf("%s:lc", redisKey)).Result()
            if err != nil {
                panic(err)
            }
            if lastCommit != lc {
                dataAttachments := []model.Attachment{
                    model.Attachment{
                        Color: "#1542e6",
                        Blocks: []model.Block{
                            m.Description,
                            m.Reviewers,
//                             m.Footer,
                        },
                    },
                }
                dataToSend, _ := json.Marshal(&model.SlackPayload{
                    Channel:        "C019QRNJ6DN",
                    ThreadTs:       ts,
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
                err = rdb.Set(ctx, fmt.Sprintf("%s:lc", redisKey), lastCommit, 0).Err()
                if err != nil {
                    panic(err)
                }
            } else {
                fmt.Println("nothing change, re-update thread")
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
                    Ts:             ts,
                    Username:       "MergeMe",
                    IconEmoji:      ":buff-mr:",
                    Attachments:    dataAttachments,
                })

                req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.update", bytes.NewBuffer(dataToSend))

                req.Header.Add("Authorization", bearer)
                req.Header.Add("Content-Type", "application/json")
                client := &http.Client{}

                resp, err := client.Do(req)
                if err != nil {
                    panic(err)
                }
                defer resp.Body.Close()

                decoder := json.NewDecoder(resp.Body)
                var t model.SlackResponsePayload
                err = decoder.Decode(&t)
                if err != nil {
                    panic(err)
                }
                defer resp.Body.Close()

                err = rdb.Set(ctx, redisKey, t.Ts, 0).Err()
                if err != nil {
                    panic(err)
                }
            }
        }
    case "merge":
        fmt.Println("merged")
        mergedBy, author := GetMergedBy(projectId, objectAttributes.Iid)
        ts, err := rdb.Get(ctx, redisKey).Result()
        if err == redis.Nil {
            dataAttachments := []model.Attachment{
                model.Attachment{
                    Color: "#1542e6",
                    Blocks: []model.Block{
                        author,
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

            decoder := json.NewDecoder(resp.Body)
            var t model.SlackResponsePayload
            err = decoder.Decode(&t)
            if err != nil {
                panic(err)
            }
            defer resp.Body.Close()

            err = rdb.Set(ctx, redisKey, t.Ts, 0).Err()
            if err != nil {
                panic(err)
            }
        } else {
            dataAttachments := []model.Attachment{
                model.Attachment{
                    Color: "#1542e6",
                    Blocks: []model.Block{
                        author,
                        m.Url,
                        mergedBy,
                        m.Footer,
                    },
                },
            }
            dataToSend, _ := json.Marshal(&model.SlackPayload{
                Channel:        "C019QRNJ6DN",
                Ts:             ts,
                Username:       "MergeMe",
                IconEmoji:      ":buff-mr:",
                Attachments:    dataAttachments,
            })

            req, err := http.NewRequest(http.MethodPost, "https://slack.com/api/chat.update", bytes.NewBuffer(dataToSend))

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
}
