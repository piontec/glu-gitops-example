#!/bin/sh

set -eo pipefail

read -p "Enter your target gitops repo URL [default: https://github.com/get-glu/gitops-example]: " repository
repository="${repository:-https://github.com/get-glu/gitops-example}"
repository="${repository%.git}"

if kubectl -n glu get secret pipeline 2>&1 >/dev/null; then
  token="$(kubectl -n glu get secret pipeline -o jsonpath='{.data.github_password}' | base64 -d)"
else
  read -s -p "Enter your GitHub personal access token [default: <empty> (read-only pipeline)]: " token
  echo ""
  echo "Creating cluster..."
fi

kind create cluster \
  --wait 120s \
  --config - <<EOF || echo "Cluster already exists"
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
name: glu-gitops-example
nodes:
- extraPortMappings:
  - containerPort: 30080 # glu pipeline
    hostPort: 30080
  - containerPort: 30081 # app in staging
    hostPort: 30081
  - containerPort: 30082 # app in production
    hostPort: 30082
  
EOF

if ! command -v go 2>&1 >/dev/null; then
  echo "Please install Go (>=1.23)"
  exit 1
fi

if ! command -v timoni 2>&1 >/dev/null; then
  go install github.com/stefanprodan/timoni/cmd/timoni@latest
fi

CONFIGURATION_REPOSITORY_URL="${repository}" \
CONFIGURATION_REPOSITORY_PASSWORD="${token}" \
  timoni bundle apply --kube-context kind-glu-gitops-example -f timoni/flux-aio.cue --runtime-from-env

echo "##########################################"
echo "#                                        #"
echo "# Pipeline Ready: http://localhost:30080 #"
echo "#                                        #"
echo "##########################################"
