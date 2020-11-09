package main

import (
    "fmt"
    _"io/ioutil"
    "log"
    "encoding/json"
    "net/http"

    gitlab  "./internal/gitlab"
    model   "./internal/model"
)

func handler(w http.ResponseWriter, r *http.Request) {
    decoder := json.NewDecoder(r.Body)
    var t model.MergeRequest
    fmt.Println(t)
    err := decoder.Decode(&t)
    if err != nil {
        return
    }

    fmt.Println(t)
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