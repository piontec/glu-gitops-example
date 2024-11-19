FROM golang:1.23

ADD . .

WORKDIR cmd/app

RUN sh -c 'go build -ldflags "-X main.SHA=$(git rev-parse HEAD)" -o /bin/app ./...'

CMD "/bin/app"
