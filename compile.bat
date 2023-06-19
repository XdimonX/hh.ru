set GOARCH=386
rem set CGO_CFLAGS=-Ic:/AutoItX
rem set CGO_LDFLAGS=-lAutoItX3_DLL
go build -ldflags="-s -w -H=windowsgui"
rem go build -ldflags="-s -w"
rem go build -ldflags="-H=windowsgui"
upx.exe --best *.exe 
