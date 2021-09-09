package model


type MergeBy struct {
    Id          int     `json:"id"`
    Name        string  `json:"name"`
    Username    string  `json:"username"`
}

//type OriAuthor struct {
    //Name        string  `json:"name"`
    //AvatarUrl   string  `json:"avatar_url"`
//}

type OriAuthor struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	State     string `json:"state"`
	AvatarURL string `json:"avatar_url"`
	WebURL    string `json:"web_url"`
}

type MergeRequestApi struct {
    Id          int         `json:"id"`
    Iid         int         `json:"iid"`
    Reviewers   []OriAuthor `json:"reviewers"`
    MergeBy     MergeBy     `json:"merged_by"`
    OriAuthor   OriAuthor   `json:"author"`
}
