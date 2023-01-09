ldflag="-s -w"
name="filesync"

mkdir release

go build -ldflags "${ldflag}" -o release/${name}_darwin_amd64 filesync.go
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "${ldflag}" -o release/${name}_linux_amd64 filesync.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags "${ldflag}" -o release/${name}_windows_amd64.exe filesync.go
CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags "${ldflag}" -o release/${name}_windows_386.exe filesync.go

upx -9 release/${name}_darwin_amd64
upx -9 release/${name}_linux_amd64
upx -9 release/${name}_windows_amd64.exe
upx -9 release/${name}_windows_386.exe
