package app

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds runtime settings.
type Config struct {
	DebugEnabled  bool
	AwsConsoleURL string
	SlackToken    string
	SlackChannel  string
}

func NewConfig() (*Config, error) {
	debugEnabled, _ := strconv.ParseBool(os.Getenv("APP_DEBUG_ENABLED"))

	cfg := Config{
		DebugEnabled: debugEnabled,
		SlackToken:   os.Getenv("APP_SLACK_TOKEN"),
		SlackChannel: os.Getenv("APP_SLACK_CHANNEL"),
	}

	missing := []string{}
	switch {
	case cfg.SlackToken == "":
		missing = append(missing, "app_slack_token")
	case cfg.SlackChannel == "":
		missing = append(missing, "app_slack_channel")
	}
	if len(missing) > 0 {
		return &Config{}, fmt.Errorf("missing env vars: %s", strings.Join(missing, ", "))
	}

	return &cfg, nil
}
