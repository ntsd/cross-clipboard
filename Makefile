run:
	go run ./cmd/main.go

bind-android:
	gomobile bind -target=android ./mobile/...

build-mobile:
	gomobile build ./mobile/...

run-mobile:
	gomobile install ./mobile/...

android-log:
	adb logcat | grep GoLog
