package events

import (
	"encoding/json"

	"github.com/slack-go/slack"
)

type BackupJobStateChange struct {
	StateChangeEvent
	BackupJobId     string `json:"backupJobId"`
	BackupVaultArn  string `json:"backupVaultArn"`
	BackupVaultName string `json:"backupVaultName"`
	ResourceArn     string `json:"resourceArn"`
	ResourceType    string `json:"resourceType"`
	State           string `json:"state"`
	StatusMessage   string `json:"statusMessage"`
	Raw             string `json:"-"`
}

func (a *BackupJobStateChange) BuildMessage() (slack.MsgOption, slack.MsgOption) {
	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Backup Job Failed", false, false))
	detailsText := "*State:*: `" + a.State + "`\n*Vault:*: `" + a.BackupVaultName + "`\n*Type:*: `" + a.ResourceType + "`\n*Resource:*:\n`" + a.ResourceArn + "`"
	details := slack.NewSectionBlock(slack.NewTextBlockObject("mrkdwn", detailsText, false, false), nil, nil)

	desc := slack.NewSectionBlock(
		slack.NewTextBlockObject("mrkdwn", a.StatusMessage, false, false),
		nil, nil,
	)
	return slack.MsgOptionText(a.StatusMessage, false), slack.MsgOptionBlocks(
		header,
		details,
		desc,
	)
}

func NewBackupJobStateChange(raw json.RawMessage) (*BackupJobStateChange, error) {
	var sce BackupJobStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &BackupJobStateChange{}, err
	}
	return &sce, nil

}
