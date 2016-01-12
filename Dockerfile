FROM golang:onbuild

ADD . /go/src/github.com/chrismar035/sudoku-generator

RUN go install github.com/chrismar035/sudoku-generator

ENTRYPOINT /go/bin/sudoku-generator
