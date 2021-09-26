package bot

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/slack-go/slack"
	"github.com/vanpt1114/mergeme/internal/model"
	"github.com/xanzy/go-gitlab"
	"strconv"
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

	diffTxt := slack.NewTextBlockObject("plain_text", "Diff", false, false)
	diffBtnEle := slack.NewButtonBlockElement("diff", "data_to_send", diffTxt)
	diffValue := model.CustomAction{
		ProjectID:      pid,
		MergeRequestID: iid,
		Clicked:        false,
	}

	if mrStatus.Approved {
		approveTxt.Text = "Approved"
		approveBtnEle.Style = "primary"
		approveValue.Clicked = true
	} else {
		approveTxt.Text = "Approve"
	}

	bApprove, err := json.Marshal(approveValue)
	if err != nil {
		panic(err)
	}
	bDiff, err := json.Marshal(diffValue)
	if err != nil {
		panic(err)
	}

	approveBtnEle.Value = string(bApprove)
	diffBtnEle.Value = string(bDiff)
	buttonBlock = slack.NewActionBlock("button", approveBtnEle, diffBtnEle)

	return buttonBlock
}

func HandleActionEvent(event slack.InteractionCallback, gl *gitlab.Client, sl *slack.Client) error {
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
				return err
			}
		case "diff":
			err := handleDiffButton(event, tmp.ProjectID, tmp.MergeRequestID, gl, sl)
			if err != nil {
				return err
			}
		default:
			fmt.Printf("%s does not supported yet\n", action.ActionID)
			return nil
		}
	}
	return nil
}


func handleApproveButton(pid, iid int, gl *gitlab.Client) error {
	_, resp, err := gl.MergeRequestApprovals.ApproveMergeRequest(pid, iid, nil)
	if err != nil {
		fmt.Println(resp.Body)
		return err
	}
	return nil
}

func handleDiffButton(event slack.InteractionCallback, pid, iid int, gl *gitlab.Client, sl *slack.Client) error {
	var sumAdded, sumRemoved int = 0, 0
	var diffStr []string
	mr, _, err := gl.MergeRequests.GetMergeRequestChanges(pid, iid, nil)
	if err != nil {
		panic(err)
	}

	changesCount, _ := strconv.Atoi(mr.ChangesCount)
	if changesCount > 5 {
		return errors.New("exceed max diff")
	}

	for _, change := range mr.Changes {
		if change.DeletedFile {
			continue
		}

		// Count numbers of lines changes
		// Should ignore MR that effects more than 100 lines
		// Because displaying 100+ lines on slack is consider spam
		added := strings.Count(change.Diff, "+  ")
		removed := strings.Count(change.Diff, "-  ")
		sumAdded = sumAdded + added
		sumRemoved = sumRemoved + removed

		// Message formatting
		diffStr = append(diffStr, "File: " + change.NewPath + "\n")
		diffStr = append(diffStr, "```" + change.Diff + "```\n")
	}

	if sumAdded > 100 || sumRemoved > 100 {
		return errors.New("exceed max diff")
	}

	diffText := slack.NewTextBlockObject("mrkdwn", strings.Join(diffStr, ""), false, false)
	diffBlock := slack.NewSectionBlock(diffText, nil, nil)

	_, _, err = sl.PostMessage(
		event.Channel.ID,
		slack.MsgOptionBlocks(diffBlock),
		slack.MsgOptionTS(event.Message.Timestamp))
	if err != nil {
		fmt.Println(err)
	}

	return nil
}
