FROM golang:1.5

ENV REPO github.com/brightbox/docker-machine-driver-brightbox

RUN go get github.com/aktau/github-release \
	github.com/brightbox/gobrightbox \
	github.com/docker/machine \
	golang.org/x/net/context \
	golang.org/x/oauth2

WORKDIR /go/src/${REPO}
ADD . /go/src/${REPO}
