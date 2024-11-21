Glu GitOps Example Pipeline
---------------------------

This repository is a demonstration of a GitOps pipeline held together with Glu.

The project consists of:

- A script for to a [Kind](https://kind.sigs.k8s.io/) Kubernetes (K8s) cluster in your local Docker.
- An example [application implemented in Go](./cmd/app) and published as a container image to [GHCR](https://docs.github.com/en/packages/working-with-a-github-packages-registry/working-with-the-container-registry).
- A promotion [pipeline implemented in Go using the Glu framework](./cmd/pipeline).
- Manifests for deploying these to your local kind cluster using [FluxCD](https://github.com/fluxcd/flux2).

All of this is triggered via a single [start script](./start.sh).
Read more below for the necessary dependencies needed to boostrap and run this example.

## Requirements

- [Docker](https://www.docker.com/)
- [Kind](https://kind.sigs.k8s.io/)
- [Go](https://go.dev/)

> The start script below will also install [Timoni](https://timoni.sh/) (using `go install` to do so).
> This is used to configure our Kind cluster.
> Big love to Timoni :heart:.

## Running

Before you get started you're going to want to do the following:

1. Fork this repo!
2. Clone your fork locally.
3. Make a note of your forks GitHub URL (likely `https://github.com/{your_username}/gitops-example`).
4. Generate a [GitHub Personal Access Token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/managing-your-personal-access-tokens) (if you want to experiment with promotions).

> You will need at-least read and write contents scope (`contents:write`).

Once you have done the above, you can run the start script found in the root of this repository.
The script will prompt you for your forks repository URL and access token (given you want to perform promotions).

```console
./start.sh
```
