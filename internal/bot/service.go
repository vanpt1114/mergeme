package bot

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/xanzy/go-gitlab"
	"strings"
)

var ctx = context.Background()

func NewGetReviewers(projectID, mrIID int, gl *gitlab.Client) (reviewers *slack.SectionBlock) {
	mr, _, err := gl.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		panic(err)
	}

	array := []string{}
	for _, val := range mr.Reviewers {
		array = append(array, fmt.Sprintf("<@%s>", val.Username))
	}

	reviewersText := slack.NewTextBlockObject("mrkdwn",
		fmt.Sprintf("*Reviewers: *" + strings.Join(array, " ")), false, false)

	reviewers = slack.NewSectionBlock(reviewersText, nil, nil)

	return reviewers
}

func NewUpdateSlackTs(r, ts string, rl *redis.Client) {
	err := rl.Set(ctx, r, ts, 0).Err()
	if err != nil {
		panic(err)
	}
}

func NewGetClosedBy(projectID, mrIID int, gl *gitlab.Client) (closedBy *slack.SectionBlock) {
	mr, _, err := gl.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		panic(err)
	}
	textBlock := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("*Closed by: *<@%s>", mr.ClosedBy.Username),
		false,
		false)
	closedBy = slack.NewSectionBlock(textBlock, nil, nil)

	return closedBy
}

func NewGetMergedBy(projectID, mrIID int, gl *gitlab.Client) (mergedBy *slack.SectionBlock) {
	mr, _, err := gl.MergeRequests.GetMergeRequest(projectID, mrIID, nil)
	if err != nil {
		panic(err)
	}
	textBlock := slack.NewTextBlockObject(
		"mrkdwn",
		fmt.Sprintf("*Merged by: *<@%s>", mr.MergedBy.Username),
		false,
		false)
	mergedBy = slack.NewSectionBlock(textBlock, nil, nil)

	return mergedBy
}

