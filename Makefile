run:
	go run ./cmd/main.go

android-install:
	gomobile install .\android\...

android-log:
	adb logcat | grep GoLog
