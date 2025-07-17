package events

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

type BackupVaultStateChange struct {
	StateChangeEvent
	BackupVaultId string `json:"backupVaultId"`
	State         string `json:"state"`
	IsLocked      string `json:"isLocked"`
	StatusMessage string `json:"statusMessage"`
	Raw           string `json:"-"`
}

func (sce *BackupVaultStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Backup Vault Change", false, false))
	blocks = append(blocks, header)

	var details []*slack.TextBlockObject
	details = append(details, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*State*\n%s", sce.State), false, false))
	details = append(details, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Locked*\n%s", sce.IsLocked), false, false))
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

func (sce *BackupVaultStateChange) IsAlertable() bool {
	return true
}

func NewBackupVaultStateChange(raw json.RawMessage) (*BackupVaultStateChange, error) {
	var sce BackupVaultStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &BackupVaultStateChange{}, err
	}
	if sce.IsLocked == "" {
		sce.IsLocked = "false"
	}
	sce.StatusMessage = "AWS Backup service backup vault was changed. Please confirm it was intended."
	return &sce, nil
}
