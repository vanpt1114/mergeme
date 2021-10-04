package main

import (
    "fmt"
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
