run:
	go run ./main.go

run-terminal:
	go run ./main.go -t

build:
	go build .

release:
	goreleaser release --rm-dist --snapshot

bind-android:
	ebitenmobile bind -target android -javapkg dev.ntsd.crossclipboard -o ./mobile/android/app/libs/cross-clipboard.aar ./mobile/.
	
android-log:
	adb logcat -c && adb logcat | grep GoLog

protogen:
	cd ./pkg/protobuf && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto
