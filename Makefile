BINARY=log-checker
MODULE=github.com/thetherington/log-checker

# These are the values we want to pass for version information
VERSION=`git describe --tags`
COMMIT=`git rev-parse HEAD`
DATE=`date +%FT%T%z`
BUILTBY=Makefile

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X ${MODULE}/cmd.Version=${VERSION} -X ${MODULE}/cmd.Commit=${COMMIT} -X ${MODULE}/cmd.Date=${DATE} -X ${MODULE}/cmd.BuiltBy=${BUILTBY}"

build:
	go build ${LDFLAGS} -o ${BINARY} main.go

docker_image:
	docker buildx build -t ${BINARY} -t ${BINARY}:${VERSION} .

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

run:
	go run main.go

.PHONY: clean build