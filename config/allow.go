package config

var AllowProject = map[int]string{
    1604: "C019QRNJ6DN",
    1052: "C01AY72GHT2",
    1714: "C01AY72GHT2", // sre/application.cluster, channel #sre
    1750: "C01AY72GHT2",
    1603: "C01FPJ7CHGQ",
    1640: "C01FPJ7CHGQ",
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
