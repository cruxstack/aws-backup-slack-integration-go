package main

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	awsevents "github.com/aws/aws-lambda-go/events"
	"github.com/joho/godotenv"

	"github.com/cruxstack/aws-backup-slack-integration-go/internal/app"
)

func main() {
	envpath := filepath.Join("..", "..", ".env")
	log.Print(envpath)
	if _, err := os.Stat(envpath); err == nil {
		_ = godotenv.Load(envpath)
	}

	cfg, err := app.NewConfig()
	if err != nil {
		log.Fatal(err)
	}
	a := app.New(cfg)

	path := filepath.Join("..", "..", "fixtures", "samples-backup-job.json")
	raw, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	var evts []awsevents.CloudWatchEvent
	if err := json.Unmarshal(raw, &evts); err != nil {
		log.Fatal(err)
	}

	for _, e := range evts {
		if err := a.Process(e); err != nil {
			log.Fatalf("process id=%s: %v", e.ID, err)
		}
	}
}
