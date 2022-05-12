run:
	go run ./cmd/main.go

android-install:
	gomobile install .\cmd\android\...

android-log:
	adb logcat | grep GoLog
