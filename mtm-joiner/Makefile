build:
	CGO_ENABLED=0 GO111MODULES=on go build -ldflags="-s -w" -o mtm-joiner .

armv8:
	CGO_ENABLED=0 GO111MODULES=on \
	GOOS=android GOARCH=arm64 \
	go build -ldflags="-s -w" -o mtm-joiner-armv8 .

.DEFAULT_GOAL := build
