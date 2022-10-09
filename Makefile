# Построить образ прежде чем собрать под windows
#docker_prebuild_image:
	#docker build -t exitstop/golang_bakend_msys2 -f docker/cross/Dockerfile .

#windows:
	##GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -v -o build/speaker.exe cmd/voice/main.go
	#./scripts/cross.sh

run: build/speaker
	./build/speaker

install_depend:
	go install github.com/playwright-community/playwright-go/cmd/playwright
	playwright install --with-deps
	#go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps

.PHONY: build/speaker
build/speaker:
	go build -v -o build/speaker cmd/speaker/main.go
	#CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -a -installsuffix cgo -v -o build/android_speaker cmd/android/main.go
