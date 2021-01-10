set GOARCH=arm
set GOOS=linux
go build -o ..\bin\pv ..\cmd\main.go ..\cmd\config.go ..\cmd\server.go ..\cmd\webservice.go

set GOARCH=386
set GOOS=windows
go build -o ..\bin\pv.exe ..\cmd\main.go ..\cmd\config.go ..\cmd\server.go ..\cmd\webservice.go