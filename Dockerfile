FROM golang:1.5

ENV REPO github.com/NeilW/docker-machine-driver-brightbox

RUN go get github.com/aktau/github-release
WORKDIR /go/src/${REPO}
ADD . /go/src/${REPO}
