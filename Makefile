BINARY=log-checker

# These are the values we want to pass for version information
VERSION=`git describe --tags`
COMMIT=`git rev-parse HEAD`
DATE=`date +%FT%T%z`
BUILTBY=Makefile

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS=-ldflags "-s -w -X github.com/thetherington/log-checker/cmd.Version=${VERSION} -X github.com/thetherington/log-checker/cmd.Commit=${COMMIT} -X github.com/thetherington/log-checker/cmd.Date=${DATE} -X github.com/thetherington/log-checker/cmd.BuiltBy=${BUILTBY}"

build:
	go build ${LDFLAGS} -o ${BINARY} main.go

# Cleans our project: deletes binaries
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi

run:
	go run main.go

.PHONY: clean build