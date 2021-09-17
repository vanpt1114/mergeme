package service

import (
	"encoding/json"
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
