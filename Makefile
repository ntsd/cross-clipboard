run:
	go run ./main.go

run-terminal:
	go run ./main.go -t

build:
	go build .

release:
	goreleaser release --snapshot --clean

bind-android:
	gomobile bind -target=android ./mobile/...

build-mobile:
	gomobile build ./mobile/...

run-mobile:
	gomobile install ./mobile/...

android-log:
	adb logcat | grep GoLog

protogen:
	cd ./pkg/protobuf && protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data.proto
