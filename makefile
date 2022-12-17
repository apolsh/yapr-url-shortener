BUILD_VERSION=1.19
BUILD_DATE=$$(date '+%d.%m.%Y')
BUILD_COMMIT=dummy

build:
	go build -ldflags "-X main.buildVersion=$(BUILD_VERSION) -X main.BuildDate=$(BUILD_DATE) -X main.BuildCommit=$(BUILD_COMMIT)" ./cmd/shortener/main.go

