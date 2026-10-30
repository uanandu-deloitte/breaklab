#!/usr/bin/env bash
set -euo pipefail

wait-for-k8s

kubectl create namespace argocd
# The applicationset CRD in this release trips kubectl apply's
# last-applied-configuration annotation size limit -- harmless here,
# this lab only needs the Application/AppProject CRDs it installs
# alongside it.
kubectl apply -n argocd -f /manifests/argocd-core-install.yaml || true

kubectl -n argocd rollout status deployment/argocd-repo-server --timeout=180s
kubectl -n argocd rollout status deployment/argocd-redis --timeout=180s
kubectl -n argocd rollout status statefulset/argocd-application-controller --timeout=180s

kubectl apply -f /manifests/project.yaml

# A real bare git repo, served by a real git daemon -- the "GitHub" for
# this lab.
mkdir -p /repos/app.git
(cd /repos/app.git && git init -q --bare)

mkdir -p /tmp/app-src/manifests/web
cp /manifests/web-deploy.yaml /tmp/app-src/manifests/web/deploy.yaml
(
  cd /tmp/app-src
  git init -q
  git config user.email "ci@example.com"
  git config user.name "ci"
  git checkout -q -b main
  git add -A
  git commit -q -m "add web manifests"
  git remote add origin /repos/app.git
  git push -q origin main
)

node_ip=$(kubectl get nodes -o jsonpath='{.items[0].status.addresses[?(@.type=="InternalIP")].address}')

# BUG: path points at manifests/app, but the repo only has manifests
# under manifests/web -- someone renamed the directory in the repo
# without updating this Application.
cat <<EOF | kubectl apply -f -
apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: web
  namespace: argocd
spec:
  project: default
  source:
    repoURL: git://${node_ip}/app.git
    targetRevision: main
    path: manifests/app
  destination:
    server: https://kubernetes.default.svc
    namespace: default
  syncPolicy:
    automated: {}
EOF
