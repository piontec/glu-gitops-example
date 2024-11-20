FROM golang:1.23

ADD . .

WORKDIR cmd/pipeline

RUN sh -c 'go build -o /bin/pipeline .'

CMD "/bin/pipeline"
