# aws-backup-integration-slack-go

**UNDER DEVELOPMENT**

AWS Lambda function that listens to **AWS Backup** events via **Amazon
EventBridge** and publishes clean, threaded messages to Slack.

## Features

* **native eventbridge trigger** – Backup events invoke the function directly
* **rich slack threads** – each finding opens a thread with severity, region,
  account and a “view in console” button
* **config-driven** – all behavior controlled by environment variables

---

## Deployment

### Prerequisites

* AWS account with AWS Backup enabled in at least one region
* Slack App with a Bot Token (`chat:write` scope) installed in your workspace
* Go ≥ 1.24
* AWS CLI configured for the deployment account

### Steps

```bash
git clone https://github.com/cruxstack/aws-backup-integration-slack.git
cd aws-backup-integration-slack

# build static Linux binary for lambda
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bootstrap

# package
zip deployment.zip bootstrap
```

## Required Environment Variables

| name                  | example                                    | purpose                                                      |
| --------------------- | ------------------------------------------ | ------------------------------------------------------------ |
| `APP_SLACK_TOKEN`     | `xoxb-…`                                   | slack bot token (store in secrets manager)                   |
| `APP_SLACK_CHANNEL`   | `C000XXXXXXX`                              | channel id to post findings                                  |
| `APP_DEBUG_ENABLED`   | `true`                                     | verbose logging & event dump                                 |

## Create Lambda Function

1. **IAM role**
   * `AWSLambdaBasicExecutionRole` managed policy
   * no additional AWS API permissions are required
2. **Lambda config**
   * Runtime: `al2023provided.al2023` (provided.al2 also works)
   * Handler: `bootstrap`
   * Architecture: `x86_64` or `arm64`
   * Upload `deployment.zip`
   * Set environment variables above
3. **EventBridge rule**
   ```json
   {
     "source": ["aws.backup"],
     "detail-type": ["Backup Job State Change"]
   }
   ```
   Target: the Lambda function.
4. **Slack App**
   * Add `chat:write` and `chat:write.public`
   * Custom bot avatar: upload AWS Backup logo in the Slack App *App Icon*
     section.


## Local Development

### Test with Samples

```bash
cp .env.example .env # edit the values
go run .
```

The sample runner replays `fixtures/samples.json` and posts to Slack exactly as
the live Lambda would.

