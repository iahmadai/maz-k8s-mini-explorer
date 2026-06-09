#!/bin/sh
set -eu

mkdir -p /home/nonroot/.kube
cp /kube/config /home/nonroot/.kube/config
sed -i 's/127.0.0.1/host.docker.internal/g' /home/nonroot/.kube/config
chmod 600 /home/nonroot/.kube/config

export KUBECONFIG=/home/nonroot/.kube/config
exec /k8s-explorer "$@"
