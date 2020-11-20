package model


type MergeBy struct {
    Id          int     `json:"id"`
    Name        string  `json:"name"`
    Username    string  `json:"username"`
}

type OriAuthor struct {
    Name        string  `json:"name"`
    AvatarUrl   string  `json:"avatar_url"`
}

type MergeRequestApi struct {
    Id          int         `json:"id"`
    Iid         int         `json:"iid"`
    MergeBy     MergeBy     `json:"merged_by"`
    OriAuthor   OriAuthor   `json:"author"`
}
