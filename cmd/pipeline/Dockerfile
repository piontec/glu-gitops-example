FROM ubuntu:24.10

RUN apt-get update && apt-get install -y \
    ca-certificates \
    git
COPY cmd/pipeline/pipeline /pipeline
COPY glu.yaml /glu.yaml

CMD ["/pipeline"]
