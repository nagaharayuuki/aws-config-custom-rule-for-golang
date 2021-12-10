GOARCH = amd64

all: clean build
clean:
	rm -rf out/

build: s3-lifecycle

s3-lifecycle:
	CGO_ENABLED=0 GOOS=linux GOARCH=${GOARCH} go build -o out/s3-lifecycle rules/s3-lifecycle/*.go
