package events

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

type BackupPlanStateChange struct {
	StateChangeEvent
	BackupPlanId  string `json:"backupPlanId"`
	VersionId     string `json:"versionId"`
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Raw           string `json:"-"`
}

func (sce *BackupPlanStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Backup Plan Change", false, false))
	blocks = append(blocks, header)

	var details []*slack.TextBlockObject
	details = append(details, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*State*\n%s", sce.State), false, false))
	blocks = append(blocks, slack.NewSectionBlock(nil, details, nil))

	if sce.StatusMessage != "" {
		desc := slack.NewSectionBlock(
			slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("```\n%s\n```", sce.StatusMessage), false, false),
			nil, nil,
		)
		blocks = append(blocks, desc)
	}
	return slack.MsgOptionText(sce.StatusMessage, false), slack.MsgOptionBlocks(blocks...)
}

func (sce *BackupPlanStateChange) IsAlertable() bool {
	return true
}

func NewBackupPlanStateChange(raw json.RawMessage) (*BackupPlanStateChange, error) {
	var sce BackupPlanStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &BackupPlanStateChange{}, err
	}
	sce.StatusMessage = "Backup plan for AWS Backup service was changed. Please confirm it was intended."
	return &sce, nil
}
