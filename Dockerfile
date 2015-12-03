FROM golang:1.5.1-alpine

RUN apk --update add git
ADD . /usr/local/go/src/github.com/NarrativeTeam/shawty

WORKDIR /usr/local/go/src/github.com/NarrativeTeam/shawty
RUN go get
RUN go build

CMD shawty
