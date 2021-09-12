package main

import (
	"encoding/json"
	"fmt"
	"github.com/slack-go/slack"
)

type Message struct {
	Author	slack.ContextElements
}

func main() {
	var author slack.Blocks
	author.BlockSet = []slack.Block{
		slack.Bloc,
	}
	//author.ContextElements = []slack.MixedElement{
	//	slack.ImageBlockElement{
	//		Type: "image",
	//		ImageURL: "aaa",
	//		AltText: "default",
	//	},
	//	slack.TextBlockObject{
	//		Type: "plain_text",
	//		Text: "hello",
	//		Emoji: true,
	//	},
	//}
	 outAuthor, _ := json.Marshal(author)
	 fmt.Println(string(outAuthor))

	var message slack.Attachment
	message.Color = "#222222"
	message.Blocks = slack.Blocks{
		BlockSet: []slack.Block{
		},
	}
	out, _ := json.Marshal(message)
	fmt.Println(string(out))

	slackClient := slack.New("xoxb-887404569732-1182478115296-mUuB9N7ZQljWaCc89XOoWoeS")
	respChannel, respTs, err := slackClient.PostMessage(
		"C015ZH0JUDC",
		slack.MsgOptionAttachments(message),
		)
	fmt.Println(respChannel, respTs)
	if err != nil {
		panic(err)
	}
	fmt.Println("Done")
}
