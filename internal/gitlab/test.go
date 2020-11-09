package gitlab

import (
    "fmt"
//     "encoding/json"

    model   "../model"
    slack   "../slack"
)

func Handle(data model.MergeRequest) {
    fmt.Println("hello")
    author := slack.Author(data.User)
    url := slack.Url(data.ObjectAttributes)
    description := slack.Description(data.ObjectAttributes)
    result := slack.Message{
        Author: author,
        Url: url,
        Description: description,
    }
    slack.PostMessage(result)
//     fmt.Println("hello")
//     fmt.Println(url)
//     author := slack.Author(data.User)
//     fmt.Println(author)
}