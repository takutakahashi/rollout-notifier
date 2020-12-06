build:
	go build -o dist/cmd cmd/cmd.go
run: build
	dist/cmd --namespace default --webhook ${WEBHOOK} --type slack
