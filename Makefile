.PHONY: run
run:
	go run cmd/workpulse/main.go

.PHONY: build.windows
build.windows:
	CC=/opt/homebrew/bin/x86_64-w64-mingw32-gcc \
	CGO_ENABLED=1 \
	GOOS=windows \
	GOARCH=amd64 \
	PKG_CONFIG_PATH=/opt/homebrew/lib/pkgconfig \
	go build -o build/workpulse.exe cmd/workpulse/main.go

.PHONY: build.linux
build.linux:
	GOOS=linux GOARCH=amd64 go build -o build/workpulse cmd/workpulse/main.go

.PHONY: build.mac
build.mac:
	GOOS=darwin GOARCH=amd64 go build -o build/workpulse cmd/workpulse/main.go

