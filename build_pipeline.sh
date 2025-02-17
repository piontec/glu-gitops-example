#!/bin/bash

go build -o cmd/pipeline/pipeline ./cmd/pipeline
docker build . -f cmd/pipeline/Dockerfile -t local/glu-pipeline:latest --load
kind load docker-image local/glu-pipeline:latest --name glu-gitops-example
kubectl -n glu rollout restart deployment pipeline
