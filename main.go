package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "bytes"
    "encoding/json"
    "net/http"

    gitlab  "github.com/vanpt1114/mergeme/internal/gitlab"
    model   "github.com/vanpt1114/mergeme/internal/model"
    slack   "github.com/vanpt1114/mergeme/internal/slack"
    config  "github.com/vanpt1114/mergeme/config/config"
)

func handler(w http.ResponseWriter, r *http.Request) {
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        panic(err)
    }
    var tmp model.MergeRequest
    err = json.Unmarshal(body, &tmp)
    if err != nil {
        panic(err)
    }
    channel := slack.CheckAllow(tmp.Project.Id)
    if channel == "" {
        fmt.Printf("ProjectID %v is not allowed\n", tmp.Project.Id)
        return
    }
//     fmt.Println(string(body))
    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))

    decoder := json.NewDecoder(r.Body)
    var t model.MergeRequest
    err = decoder.Decode(&t)
    if err != nil {
        return
    }
    gitlab.Handle(t)
}

func health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}


func main() {
    http.HandleFunc("/health/live", health)
    http.HandleFunc("/health/ready", health)
    http.HandleFunc("/merge-me", handler)
    fmt.Println("MergeMe is running on at http://127.0.0.1:10080/merge-me")
    log.Fatal(http.ListenAndServe(":10080", nil))
}