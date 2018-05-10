
PACKAGE=github.com/troykinsella/qc
BINARY=qc
VERSION=0.0.1

LDFLAGS=-ldflags "-X main.AppVersion=${VERSION}"

build:
	go build ${LDFLAGS} ${PACKAGE}

install:
	go install ${LDFLAGS}

test:
	go test ${PACKAGE}/...

coverage:
	go test -cover ${PACKAGE}/...

dist:
	gox ${LDFLAGS} \
		-arch="amd64" \
		-os="darwin linux windows" \
		-output="${BINARY}_{{.OS}}_{{.Arch}}" \
		${PACKAGE}

clean:
	test -f ${BINARY} && rm ${BINARY} || true
	rm ${BINARY}_* || true

.PHONY: build install test coverage dist release clean
