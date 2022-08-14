# Построить образ прежде чем собрать под windows
#docker_prebuild_image:
	#docker build -t exitstop/golang_bakend_msys2 -f docker/cross/Dockerfile .

#windows:
	##GOOS=windows GOARCH=amd64 CGO_ENABLED=1 $(GOBUILD) -v -o build/speaker.exe cmd/voice/main.go
	#./scripts/cross.sh

android:
	go run cmd/android/main.go -ip 192.168.0.177

android_work:
	go run cmd/speaker/main.go -ip 192.168.88.20

google_speech: build/speaker
	./build/speaker -google_speech true

read_ru:
	go run cmd/speaker/ru.go

install_depend:
	go install github.com/playwright-community/playwright-go/cmd/playwright
	playwright install --with-deps
	#go run github.com/playwright-community/playwright-go/cmd/playwright install --with-deps

.PHONY: build/speaker
build/speaker:
	go build -v -o build/speaker cmd/speaker/main.go
	#CGO_ENABLED=0 go build -ldflags '-w -extldflags "-static"' -a -installsuffix cgo -v -o build/android_speaker cmd/android/main.go
