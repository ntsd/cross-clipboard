run:
	go run ./cmd/main.go

run-ui:
	go run ./ui/main.go

android-install:
	cd ui && fyne install -os android -appID dev.ntsd.cross.clipboard

android-build:
	cd ui && fyne package -os android -appID dev.ntsd.cross.clipboard

android-log:
	adb logcat | grep GoLog
