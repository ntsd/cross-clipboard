run:
	go run ./cmd/cross-clipboard.go

build:
	go build ./cmd/cross-clipboard.go

bind-android:
	gomobile bind -target=android ./mobile/...

build-mobile:
	gomobile build ./mobile/...

run-mobile:
	gomobile install ./mobile/...

android-log:
	adb logcat | grep GoLog
