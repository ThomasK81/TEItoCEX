compile:
	echo "Compiling for every OS and Platform"
	GOOS=freebsd GOARCH=386 go build -o bin/TEItoCEX-FreeBDS-386 CTSExtract.go
	GOOS=darwin GOARCH=amd64 go build -o bin/TEItoCEX-OSX CTSExtract.go
	GOOS=linux GOARCH=386 go build -o bin/TEItoCEX-Linux-386 CTSExtract.go
	GOOS=windows GOARCH=386 go build -o bin/TEItoCEX-Windows-386 CTSExtract.go
