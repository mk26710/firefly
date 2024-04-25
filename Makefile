PROJECT_NAME=firefly

build_linux:
	GOARCH=amd64 GOOS=linux   go build -o .output/${PROJECT_NAME}-amd64     -ldflags "-s -w" cmd/${PROJECT_NAME}/main.go
	GOARCH=arm64 GOOS=linux   go build -o .output/${PROJECT_NAME}-arm64     -ldflags "-s -w" cmd/${PROJECT_NAME}/main.go

build_win:
	GOARCH=amd64 GOOS=windows go build -o .output/${PROJECT_NAME}-amd64.exe -ldflags "-s -w" cmd/${PROJECT_NAME}/main.go
	GOARCH=arm64 GOOS=windows go build -o .output/${PROJECT_NAME}-arm64.exe -ldflags "-s -w" cmd/${PROJECT_NAME}/main.go

build: build_linux build_win
