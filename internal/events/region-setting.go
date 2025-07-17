package events

import (
	"encoding/json"
	"fmt"

	"github.com/slack-go/slack"
)

type RegionSettingStateChange struct {
	StateChangeEvent
	State         string `json:"state"`
	StatusMessage string `json:"statusMessage"`
	Raw           string `json:"-"`
}

func (sce *RegionSettingStateChange) SlackMessage() (slack.MsgOption, slack.MsgOption) {
	var blocks []slack.Block

	header := slack.NewHeaderBlock(slack.NewTextBlockObject("plain_text", "Region Setting Change", false, false))
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

func (sce *RegionSettingStateChange) IsAlertable() bool {
	return true
}

func NewRegionSettingStateChange(raw json.RawMessage) (*RegionSettingStateChange, error) {
	var sce RegionSettingStateChange
	if err := json.Unmarshal(raw, &sce); err != nil {
		return &RegionSettingStateChange{}, err
	}
	sce.StatusMessage = "AWS Backup service region settings changed. Please confirm it was intended."
	return &sce, nil
}
