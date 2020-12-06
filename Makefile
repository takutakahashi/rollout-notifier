build:
	go build -o dist/cmd cmd/cmd.go
run: build
	dist/cmd --namespace default --config .ignore/config.yaml --type noop
run_slack: build
	dist/cmd --namespace default --config .ignore/config.yaml --type slack
