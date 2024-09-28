BUILDDIR := ./build
APPNAME := nautilus
SERVERNAME := ${APPNAME}-server
CLIENTNAME := ${APPNAME}-client
HOST := localhost
PORT := 3030

build-server:
	@go build -o ${BUILDDIR}/${SERVERNAME} ./cmd/server

build-client:
	@go build -o ${BUILDDIR}/${CLIENTNAME} ./cmd/client

build: build-server build-client

run-server: build-server
	@${BUILDDIR}/${SERVERNAME} -host ${HOST} -port ${PORT}

run-client: build-client
	@${BUILDDIR}/${CLIENTNAME} -host ${HOST} -port ${PORT}
