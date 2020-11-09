package slack

import (
//     "fmt"
    "bytes"
    "encoding/json"
    "net/http"
)

type Message struct {
    Author      Block
    Url         Block
    Description Block
}

func PostMessage(m Message) {
    url := "https://slack.com/api/chat.postMessage"
    bearer := "Bearer " + "xxxx"

    dataAttachments := []Attachment{
        Attachment{
            Color: "#1542e6",
            Blocks: []Block{
                m.Author,
                m.Url,
                m.Description,
            },
        },
    }

    dataToSend, _ := json.Marshal(SlackPayload{
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
}
