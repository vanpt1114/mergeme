package main

import (
    "fmt"
    "github.com/vanpt1114/mergeme/config"
    "github.com/vanpt1114/mergeme/internal/service"
    "log"
    "net/http"
)

var cfg *config.Config

//func handler(w http.ResponseWriter, r *http.Request) {
//    body, err := ioutil.ReadAll(r.Body)
//    defer r.Body.Close()
//    if err != nil {
//        panic(err)
//    }
//    var tmp gitlab.MergeEvent
//    err = json.Unmarshal(body, &tmp)
//    if err != nil {
//        panic(err)
//    }
//    channel := config.CheckAllow(tmp.Project.ID)
//    if channel == "" {
//       fmt.Printf("ProjectID %v is not allowed\n", tmp.Project.ID)
//       return
//    }
//    if tmp.ObjectAttributes.WorkInProgress == true {
//       fmt.Println("WIP is not allow")
//       return
//    }
//    r.Body = ioutil.NopCloser(bytes.NewBuffer(body))
//    decoder := json.NewDecoder(r.Body)
//    var t gitlab.MergeEvent
//    err = decoder.Decode(&t)
//    if err != nil {
//       return
//    }
//    svc := service.NewService(cfg)
//    svc.HandleEvent(t)
//}


func health(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
}


func main() {
    cfg = config.Load()
    svc := service.NewService(cfg)
    http.HandleFunc("/health/live", health)
    http.HandleFunc("/health/ready", health)
    http.HandleFunc("/merge-me", svc.HandleRequest)
    fmt.Println("MergeMe is running on at http://0.0.0.0:10080/merge-me")
    log.Fatal(http.ListenAndServe(":10080", nil))
}
