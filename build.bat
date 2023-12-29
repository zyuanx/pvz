go build -ldflags "-s -w -H=windowsgui"

upx -9 pvz.exe