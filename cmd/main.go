package main

import (
    "encoding/json"
    "fmt"
    "github.com/slack-go/slack"
    "github.com/vanpt1114/mergeme/config"
    "github.com/vanpt1114/mergeme/internal/service"
    "log"
    "net/http"
    "os"
)

var cfg *config.Config

func health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}

type CustomAction struct {
    ProjectID   int `json:"project_id"`
    MrID        int `json:"iid"`
}

func approve(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    //fmt.Println(r.Header.Get("Content-Type"))
    r.ParseForm()
    data := r.PostFormValue("payload")
    //fmt.Println(data)
    var t slack.InteractionCallback
    err := json.Unmarshal([]byte(data), &t)
    if err != nil {
        panic(err)
    }
    dataJson, _ := json.Marshal(t)
    fmt.Println(string(dataJson))

    for _, action := range t.ActionCallback.BlockActions {
        if action.ActionID == "approve" {
            var tmp CustomAction
            err := json.Unmarshal([]byte(action.Value), &tmp)
            if err != nil {
                panic(err)
            }
            fmt.Println(tmp.ProjectID)
        }
    }
}

func main() {
    cfg = config.Load()
    svc := service.NewService(cfg)
    http.HandleFunc("/health/live", health)
    http.HandleFunc("/health/ready", health)
    http.HandleFunc("/merge-me", svc.HandleRequest)
    http.HandleFunc("/merge-me/approve", svc.HandleAction)
    fmt.Println("MergeMe is running on at http://0.0.0.0:10080/merge-me")
    log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), nil))
}
