package service

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/bot"
	"github.com/xanzy/go-gitlab"
	"io/ioutil"
	"net/http"
)

func (s *Service) HandleRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		panic(err)
	}

	var requestBody = body
	var event gitlab.MergeEvent
	err = json.Unmarshal(requestBody, &event)
	if err != nil {
		panic(err)
	}

	if s.shouldSkipMergeRequest(event) {
		w.WriteHeader(http.StatusBadRequest)
		_, err = w.Write([]byte("Skip MR"))
		if err != nil {
			panic(err)
		}
		return
	}

	// HandleEvent event
	s.HandleEvent(event)
}

func (s *Service) HandleAction(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != "application/x-www-form-urlencoded" {
		return
	}

	if err := r.ParseForm(); err != nil {
		panic(err)
	}

	requestBody := r.PostFormValue("payload")
	var event slack.InteractionCallback
	err := json.Unmarshal([]byte(requestBody), &event)
	if err != nil {
		panic(err)
	}

	//if s.shouldSkipAction(event) {
	//	fmt.Println("button has been clicked")
	//	return
	//}

	err = bot.HandleActionEvent(event, s.gitlab, s.slack)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, err := w.Write([]byte("error when clicking button"))
		if err != nil {
			fmt.Println(err)
		}
	}

	w.WriteHeader(http.StatusOK)
}
