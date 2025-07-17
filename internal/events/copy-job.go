package events

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/slack-go/slack"
)

type CopyJobStateChange struct {
	StateChangeEvent
	CopyJobId     string `json:"copyJobId"`
	ResourceArn   string `json:"resourceArn"`
	ResourceType  string `json:"resourceType"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Raw           string `json:"-"`
}

func (sce *CopyJobStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Copy Job Failed", false, false))
	blocks = append(blocks, header)

	var detailFields []*slack.TextBlockObject
	detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*State*\n%s", sce.State), false, false))

	if sce.ResourceArn != "" {
		rArn, rArnErr := arn.Parse(sce.ResourceArn)
		if rArnErr == nil {
			detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Account*\n%s", rArn.AccountID), false, false))
			detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Type*\n%s", sce.ResourceType), false, false))
			detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Region*\n%s", rArn.Region), false, false))
			detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Resource*\n%s", rArn.Resource), false, false))
		}
	}
	details := slack.NewSectionBlock(nil, detailFields, nil)
	blocks = append(blocks, details)

	if sce.StatusMessage != "" {
		desc := slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", sce.StatusMessage), false, false),
			nil, nil,
		)
		blocks = append(blocks, desc)
	}

	return slack.MsgOptionText(sce.StatusMessage, false), slack.MsgOptionBlocks(blocks...)
}

func (sce *CopyJobStateChange) IsAlertable() bool {
	failedStates := []string{"FAILED"}
	return slices.Contains(failedStates, sce.State)
}

func NewCopyJobStateChange(raw json.RawMessage) (*CopyJobStateChange, error) {
	var sce CopyJobStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &CopyJobStateChange{}, err
	}
	sce.StatusMessage = strings.Trim(sce.StatusMessage, "\"")
	return &sce, nil
}
