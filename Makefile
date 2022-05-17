run:
	go run ./cmd/main.go

run-ui:
	go run ./ui/main.go

android-log:
	adb logcat | grep GoLog
