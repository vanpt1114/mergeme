package config

var AllowProject = map[int]string{
    1604: "C019QRNJ6DN",
}

func CheckAllow(projectId int) (channel string) {
    channel, ok := AllowProject[projectId]
    if ok == false {
        return ""
    }
    return
}
