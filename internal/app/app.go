package app

import (
	"fmt"

	awsEvent "github.com/aws/aws-lambda-go/events"
	"github.com/cruxstack/aws-backup-slack-integration-go/internal/events"
	"github.com/slack-go/slack"
)

// App orchestrates parsing and Slack publishing.
type App struct {
	Config      *Config
	SlackClient *slack.Client
}

// New returns an initialized App.
func New(cfg *Config) *App {
	return &App{
		Config:      cfg,
		SlackClient: slack.New(cfg.SlackToken),
	}
}

// ParseEvent converts raw CloudWatch event into a typed structure.
func (a *App) ParseEvent(e awsEvent.CloudWatchEvent) (events.StateChangeEvent, error) {
	switch e.DetailType {
	case "Backup Job State Change":
		return events.NewBackupJobStateChange(e.Detail)
	case "Copy Job State Change":
		return events.NewCopyJobStateChange(e.Detail)
	case "Restore Job State Change":
		return events.NewRestoreJobStateChange(e.Detail)
	default:
		return nil, fmt.Errorf("unknown cloudwatch event type: %s", e.DetailType)
	}
}

// Process handles a single CloudWatch event end-to-end.
func (a *App) Process(evt awsEvent.CloudWatchEvent) error {
	sce, err := a.ParseEvent(evt)
	if err != nil || !sce.IsAlertable() {
		return err
	}
	m0, m1 := sce.SlackMessage()
	_, _, err = a.SlackClient.PostMessage(a.Config.SlackChannel, m0, m1)
	return err
}
