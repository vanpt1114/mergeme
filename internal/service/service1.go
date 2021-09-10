package service

import (
    "encoding/json"
    "github.com/vanpt1114/mergeme/internal/model"
    "net/http"
)


func UpdateSlackTs(r, ts string) {
    err := rdb.Set(ctx, r, ts, 0).Err()
    if err != nil {
        panic(err)
    }
}

func DecodeSlackResponse(r *http.Response) string {
    decoder := json.NewDecoder(r.Body)
    var t model.SlackResponsePayload
    err := decoder.Decode(&t)
    if err != nil {
        panic(err)
    }
    if t.Ok != true {
        panic(t.Error)
    }
    return t.Ts
}

