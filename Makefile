run:
	go run ./cmd/main.go

mobile-install:
	gomobile install .\mobile\...

android-log:
	adb logcat | grep GoLog
