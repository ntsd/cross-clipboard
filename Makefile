run:
	go run ./cmd/main.go

build-mobile:
	gomobile build ./mobile/...

run-mobile:
	gomobile install ./mobile/...

android-log:
	adb logcat | grep GoLog
