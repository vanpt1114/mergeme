package bot

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/model"
	"github.com/xanzy/go-gitlab"
	"strings"
)

var ctx = context.Background()

func GetReviewers(projectID, mrIID int, gl *gitlab.Client) (reviewers *slack.SectionBlock) {
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

func UpdateSlackTs(r, ts string, rl *redis.Client) {
	err := rl.Set(ctx, r, ts, 0).Err()
	if err != nil {
		panic(err)
	}
}

func GetClosedBy(projectID, mrIID int, gl *gitlab.Client) (closedBy *slack.SectionBlock) {
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

func GetMergedBy(projectID, mrIID int, gl *gitlab.Client) (mergedBy *slack.SectionBlock) {
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

func ActionButton(pid int, iid int, gl *gitlab.Client) (buttonBlock *slack.ActionBlock) {
	mrStatus, _, err := gl.MergeRequestApprovals.GetConfiguration(pid, iid)
	if err != nil {
		panic(err)
	}

	approveTxt := slack.NewTextBlockObject("plain_text", "", false, false)
	approveBtnEle := slack.NewButtonBlockElement("approve", "data_to_send", approveTxt)
	approveValue := model.CustomAction{
		ProjectID:      pid,
		MergeRequestID: iid,
		Clicked:        false,
	}

	diffTxt := slack.NewTextBlockObject("plain_text", "", false, false)
	diffBtnEle := slack.NewButtonBlockElement("diff", "data_to_send", diffTxt)

	if mrStatus.Approved {
		approveTxt.Text = "Approved"
		approveBtnEle.Style = "primary"
		approveValue.Clicked = true
	} else {
		approveTxt.Text = "Approve"
	}

	b, err := json.Marshal(approveValue)
	if err != nil {
		panic(err)
	}
	approveBtnEle.Value = string(b)
	buttonBlock = slack.NewActionBlock("button", approveBtnEle, diffBtnEle)

	return buttonBlock
}

func HandleActionEvent(event slack.InteractionCallback, gl *gitlab.Client) {
	for _, action := range event.ActionCallback.BlockActions {
		var tmp model.CustomAction
		err := json.Unmarshal([]byte(action.Value), &tmp)
		if err != nil {
			panic(err)
		}

		// Only non-clicked button can go to this function
		// because we have check for clicked button before that
		switch action.ActionID {
		case "approve":
			err := handleApproveButton(tmp.ProjectID, tmp.MergeRequestID, gl)
			if err != nil {
				fmt.Println(err)
			}
		default:
			fmt.Printf("%s does not supported yet\n", action.ActionID)
		}
	}
}


func handleApproveButton(pid, iid int, gl *gitlab.Client) error {
	_, resp, err := gl.MergeRequestApprovals.ApproveMergeRequest(pid, iid, nil)
	if err != nil {
		fmt.Println(resp.Body)
	}
	return nil
}