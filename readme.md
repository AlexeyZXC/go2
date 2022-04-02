Process to build in windows under linux+arm platform:
1. set env:   go env -w GOARCH=arm GOOS=linux
2. build: go build main.go