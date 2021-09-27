package config

import (
    "errors"
    "github.com/vanpt1114/mergeme/internal/model"
    "gopkg.in/yaml.v2"
    "io/ioutil"
    "os"
)

var AllowProject = make(map[int]string)

func CheckAllow(projectId int) (channel string, err error) {
    channel, exist := AllowProject[projectId]
    if !exist {
        return "", errors.New("project is not allowed")
    }
    return channel, nil
}

func InitProjectsMapping() error {
    var data model.ProjectsMapping

    pwd, _ := os.Getwd()
    yamlFile, err := ioutil.ReadFile(pwd + "/default.yaml")
    if err != nil {
        panic(err)
    }

    err = yaml.Unmarshal(yamlFile, &data)
    if err != nil {
        panic(err)
    }

    for _, value := range data.Projects {
        AllowProject[value.ID] = value.Channel
    }

    return nil
}