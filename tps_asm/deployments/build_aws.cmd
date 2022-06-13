set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
mkdir .\dist
go build -o .\dist\main .\cmd\awslambda\main.go
".\3rd party\zip\7z.exe" a .\dist\main.zip .\dist\main
