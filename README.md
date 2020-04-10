![Alerty Logo](https://storage.googleapis.com/alerty-public/landing/logo.png)


# Alerty

Alerty is a Golang project for monitoring and alerting when a Website or Socket is down.

It uses GCP (Cloudfuctions and Pub/Sub).

_This project is for fun so... enjoy it._ :smile:

https://alerty.online

## Requirements

* Go: 1.13

* Mongo Database: `docker-compose -f deployments/docker-compose.yml up -d mongo`

* Configure environments: `cp .env.template .env`

## Start Cloud Functions

```bash
# Send Email:
gcloud functions deploy SendEmail --region us-east4 --runtime go111 --trigger-topic send_email --set-env-vars mailgun_key=<mailgunkey>

# Send SMS:
gcloud functions deploy SendSMS --region us-east4 --runtime go111 --trigger-topic send_sms --set-env-vars sid=<twilio_sid> --set-env-vars token=<twilio_token>

# Messenger:
gcloud functions deploy Messenger --region us-east4 --runtime go111 --trigger-topic messenger --set-env-vars project_id=<gcp_project_id>

# CheckURL:
gcloud functions deploy CheckUrl --region us-east4 --runtime go111 --trigger-topic check-url --region us-east4 --set-env-vars TOKEN=<BrainToken>

# Send Request:
gcloud functions deploy SendRequest --region us-east4 --runtime go111 --trigger-topic test-request

# Send Slack Message:
gcloud functions deploy SendSlackMessage --region us-east4 --runtime go111 --trigger-topic send-slack-message

# Robot RUN:
gcloud functions deploy RobotRun --region us-east4 --runtime go111 --trigger-topic robot-run --region us-east4 --set-env-vars TOKEN=<BrainToken>

# Cloud Function Check:
gcloud functions logs read --limit 50
```

## Testing process

```bash
# Run populate script:
go run cmd/populate/create.go

# Run Brain API:
go run cmd/brain/brain.go

# Expose via Ngrok and get BrainURL:
ngrok http 3000

# Run Controller:
go run cmd/controller/controller.go

## Runners

# Website Runner:
go run cmd/websites-cron/run.go

# Socket Runner:
go run cmd/sockets-cron/run.go

# Robot Runner:
go run cmd/robots-cron/run.go

# Check Cloud Function:
gcloud functions logs read --limit 30 --region us-east4
```

## Contributing
Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.

Please make sure to update tests as appropriate.

## Contributors
* [Juan P. Arias](https://github.com/kurkop)

* [Jhon Ramirez](https://github.com/jhoniscoding)

## License

[LGPL-3](https://www.gnu.org/licenses/lgpl-3.0.en.html)
Copyright (c) 2020 [Kurlabs](http://kurlabs.com)