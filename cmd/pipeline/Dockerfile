FROM golang:1.23

ADD . .

WORKDIR cmd/app

RUN sh -c 'go build -o /bin/pipeline ./cmd/pipeline/...'

CMD "/bin/pipeline"
