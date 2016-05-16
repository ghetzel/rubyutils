all: vendor fmt test

update:
	glide up --strip-vcs --update-vendored

vendor:
	go list github.com/Masterminds/glide
	glide install --strip-vcs --update-vendored

fmt:
	gofmt -w .
	gofmt -w encoding/..

test:
	go test -v ./encoding/...
