all: build

fmt:
	gofmt -l -w -s */

build: fmt 
	cd cmd/alarm && go build -v

install: fmt
	cd cmd/alarm && go install

clean:
	cd cmd/alarm && go clean