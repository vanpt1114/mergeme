package model


type SlackResponsePayload struct {
    Ok  bool    `json:"ok"`
    Ts  string  `json:"ts"`
}


type Child struct {
    Type     string  `json:"type,omitempty"`
    ImageUrl string  `json:"image_url,omitempty"`
    AltText  string  `json:"alt_text,omitempty"`
    Text     string  `json:"text,omitempty"`
    Emoji    bool    `json:"emoji,omitempty"`
}


type Block struct {
    Type        string      `json:"type"`
    Elements    *[]Child    `json:"elements,omitempty"`
    Text        *Child      `json:"text,omitempty"`
}

type Attachment struct {
    Color   string  `json:"color"`
    Blocks  []Block `json:"blocks"`
}

type SlackPayload struct {
    Channel         string          `json:"channel"`
    ThreadTs        string          `json:"thread_ts"`
    Ts              string          `json:"ts"`
    Username        string          `json:"username"`
    IconEmoji       string          `json:"icon_emoji"`
    Attachments     []Attachment    `json:"attachments"`
}
