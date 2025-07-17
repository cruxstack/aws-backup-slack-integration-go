package main

import (
	"context"
	"encoding/json"
	"log"
	"sync"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/cruxstack/aws-backup-slack-integration-go/internal/app"
)

/*
   Lambda wrapper helpers
*/

var (
	once    sync.Once
	a       *app.App
	initErr error
)

func LambdaHandler(_ context.Context, evt awsevents.CloudWatchEvent) error {
	once.Do(func() {
		cfg, err := app.NewConfig()
		if err != nil {
			initErr = err
			return
		}
		a = app.New(cfg)
	})

	if initErr != nil {
		// returning (not panic) preserves clean CW metrics
		return initErr
	}

	if a.Config.DebugEnabled {
		j, _ := json.Marshal(evt)
		log.Print(string(j))
	}

	return a.Process(evt)
}

func main() {
	lambda.Start(LambdaHandler)
}
