all:
	@make build
	@make archive
	@make upload

archive:
	zip calendar_notifier.zip notifier credentials.json token.json

upload:
	# aws s3 cp calendar_notifier.zip s3://lambda-okiba/
	:

build:
	GOOS=linux GOARCH=amd64 go build -o notifier notifier.go client.go cal.go plan.go
