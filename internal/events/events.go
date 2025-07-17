package events

import (
	"github.com/slack-go/slack"
)

type StateChangeEvent interface {
	BuildMessage() (slack.MsgOption, slack.MsgOption)
}
