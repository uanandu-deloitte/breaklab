# k8sd-007: The Cluster That Stopped Listening to Git

Real k3s, real ArgoCD (core components only: application-controller,
repo-server, redis -- no UI/API server needed), and a real git daemon
serving a real bare repo. The `web` Application's `source.path` points
at a directory that no longer exists in the repo after a rename,
leaving it permanently out of sync.
