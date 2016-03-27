FROM golang:1.6-alpine

COPY exercisers /go/src/github.com/calavera/docker-flames/
RUN cd /go/src/github.com/calavera/docker-flames/ && go install ./...
