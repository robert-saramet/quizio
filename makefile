run:
	go run .
linux:
	GOOS=linux GOARCH=amd64 go build -o bin/quizio-amd64-linux .
	cp *yaml bin/
windows:
	GOOS=windows GOARCH=amd64 go build -o bin/quizio-amd64.exe .