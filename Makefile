BUILDDIR := ./build
APPNAME := nautilus
SERVERNAME := ${APPNAME}-server
CLIENTNAME := ${APPNAME}-client

build-server:
	@go build -o ${BUILDDIR}/${SERVERNAME} ./cmd/server

build-client:
	@go build -o ${BUILDDIR}/${CLIENTNAME} ./cmd/client

build: build-server build-client

run-server: build-server
	@${BUILDDIR}/${SERVERNAME}

run-client: build-client
	@${BUILDDIR}/${CLIENTNAME}
