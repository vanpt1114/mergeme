package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/vanpt1114/mergeme/config"
    "github.com/xanzy/go-gitlab"
    "io/ioutil"
    "log"
    "net/http"
    "os"

    "github.com/vanpt1114/mergeme/internal/service"
)

func handler(w http.ResponseWriter, r *http.Request) {
    GITLAB_TOKEN := os.Getenv("GITLAB_TOKEN")
    GITLAB_URL := os.Getenv("GITLAB_URL")
    SLACK_TOKEN := os.Getenv("SLACK_TOKEN")
    body, err := ioutil.ReadAll(r.Body)
    defer r.Body.Close()
    if err != nil {
        panic(err)
    }
    var tmp gitlab.MergeEvent
    err = json.Unmarshal(body, &tmp)
    if err != nil {
        panic(err)
    }
    channel := config.CheckAllow(tmp.Project.ID)
    if channel == "" {
       fmt.Printf("ProjectID %v is not allowed\n", tmp.Project.ID)
       return
    }
    if tmp.ObjectAttributes.WorkInProgress == true {
       fmt.Println("WIP is not allow")
       return
    }
    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
    decoder := json.NewDecoder(r.Body)
    var t gitlab.MergeEvent
    err = decoder.Decode(&t)
    if err != nil {
       return
    }
    svc := service.NewService(GITLAB_TOKEN, GITLAB_URL, SLACK_TOKEN)
    svc.Handle(t)
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
