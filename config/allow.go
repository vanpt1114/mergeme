package config

var AllowProject = map[int]string{
    1038: "C019QRNJ6DN", // workloads develop, channel #new-test-alerts
    1067: "C019QRNJ6DN", // workloads staging, channel #new-test-alerts
    1068: "C019QRNJ6DN", // workloads production, channel #new-test-alerts
}

func CheckAllow(projectId int) (channel string) {
    channel, ok := AllowProject[projectId]
    if ok == false {
        return ""
    }
    return
}
