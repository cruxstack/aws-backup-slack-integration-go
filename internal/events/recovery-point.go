package events

import (
	"encoding/json"
	"fmt"
	"slices"

	"github.com/aws/aws-sdk-go-v2/aws/arn"
	"github.com/slack-go/slack"
)

type RecoveryPointStateChange struct {
	StateChangeEvent
	BackupVaultArn  string `json:"backupVaultArn"`
	BackupVaultName string `json:"backupVaultName"`
	ResourceArn     string `json:"resourceArn"`
	ResourceType    string `json:"resourceType"`
	State           string `json:"state"`
	Status          string `json:"status"`
	StatusMessage   string `json:"statusMessage"`
	DeletedBy       string `json:"deletedBy"`
	Raw             string `json:"-"`
}

func (sce *RecoveryPointStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Recovery Point Changed", false, false))
	blocks = append(blocks, header)

	var detailFields []*slack.TextBlockObject
	detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*State*\n%s", sce.State), false, false))

	if sce.BackupVaultName != "" {
		detailFields = append(detailFields, slack.NewTextBlockObject("mrkdwn", fmt.Sprintf("*Vault*\n%s", sce.BackupVaultName), false, false))
	}
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

func (sce *RecoveryPointStateChange) IsAlertable() bool {
	failedStates := []string{"DELETED"}
	return slices.Contains(failedStates, sce.State) && sce.DeletedBy == "MANUAL_DELETE"
}

func NewRecoveryPointStateChange(raw json.RawMessage) (*RecoveryPointStateChange, error) {
	var sce RecoveryPointStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &RecoveryPointStateChange{}, err
	}
	if sce.State == "" {
		sce.State = sce.Status
	}
	sce.StatusMessage = "Recovery point (aka backup) changed within the AWS Backup service. Please confirm it was intended."
	return &sce, nil
}
