FROM golang:1.23

ADD . /src

WORKDIR /src

RUN sh -c 'go build -o /bin/pipeline ./cmd/pipeline'

CMD "/bin/pipeline"
