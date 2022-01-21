APP=toil

env GOOS=linux GOARCH=amd64 go build -o $APP-linux-amd64 main.go
env GOOS=darwin GOARCH=amd64 go build -o $APP-darwin-amd64 main.go
