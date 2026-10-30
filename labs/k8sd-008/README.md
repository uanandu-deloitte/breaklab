# k8sd-008: Helm Rendered the Wrong Port

Real k3s, real `helm` (v4). The chart's Service hardcodes `targetPort:
8080`, but the container actually listens on port 80 — the Deployment
is healthy, but nothing reaches it through the Service.
