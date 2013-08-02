all: subzer test

get-deps:
	go get "github.com/jessevdk/go-flags"
	go get "github.com/kdar/stringio"

cross:
	GOARCH=386   GOOS=linux   go build -o subzer_linux32 subzer.go
	GOARCH=amd64 GOOS=linux   go build -o subzer_linux64 subzer.go
	GOARCH=amd64 GOOS=darwin  go build -o subzer_osx64   subzer.go
	GOARCH=386   GOOS=windows go build -o subzer32.exe   subzer.go
	GOARCH=amd64 GOOS=windows go build -o subzer64.exe   subzer.go

subzer: subzer.go
	go build subzer.go

format:
	gofmt -s -tabs=false -tabwidth=4 -w=true subzer.go
	gofmt -s -tabs=false -tabwidth=4 -w=true subzer_test.go

test: subzer_test.go subzer.go
	go test -v

clean:
	rm -f subzer subzer.test subzer_linux32 subzer_linux64 subzer_osx64 subzer32.exe subzer64.exe
