package events

import (
	"encoding/json"
	"fmt"
	"slices"
	"strings"

	"github.com/slack-go/slack"
)

type RestoreJobStateChange struct {
	StateChangeEvent
	RestoreJobId  string `json:"restoreJobId"`
	ResourceType  string `json:"resourceType"`
	State         string `json:"state"`
	Status        string `json:"status"`
	StatusMessage string `json:"statusMessage"`
	Raw           string `json:"-"`
}

func (sce *RestoreJobStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Restore Job Failed", false, false))
	blocks = append(blocks, header)

	var detailFields []*slack.TextBlockObject
	detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*State*\n%s", sce.State), false, false))

	if sce.ResourceType != "" {
		detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Type*\n%s", sce.ResourceType), false, false))
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

func (sce *RestoreJobStateChange) IsAlertable() bool {
	failedStates := []string{"FAILED"}
	return slices.Contains(failedStates, sce.Status)
}

func NewRestoreJobStateChange(raw json.RawMessage) (*RestoreJobStateChange, error) {
	var sce RestoreJobStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &RestoreJobStateChange{}, err
	}
	if sce.State == "" {
		sce.State = sce.Status
	}
	sce.StatusMessage = strings.Trim(sce.StatusMessage, "\"")
	return &sce, nil
}
