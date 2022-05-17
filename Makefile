run:
	go run ./cmd/main.go

run-ui:
	go run ./ui/main.go

bundle-assets:
	fyne bundle --package assets --prefix Resource ./assets/*.png > ./assets/bundled.go

android-install:
	fyne install -os android -appID dev.ntsd.cross.clipboard

android-build:
	fyne package -os android -appID dev.ntsd.cross.clipboard

android-log:
	adb logcat | grep GoLog
