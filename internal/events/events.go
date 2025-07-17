package events

import (
	"github.com/slack-go/slack"
)

type StateChangeEvent interface {
	IsAlertable() bool
	SlackMessage() (slack.MsgOption, slack.MsgOption)
}
