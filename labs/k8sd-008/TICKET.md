# Helm Rendered the Wrong Port

The `web` Helm release looks fine — `helm status web` is deployed, the
Deployment is Ready, the Service exists. Nothing reaches it. Every
request to the Service's port just hangs or refuses.

## Tasks

1. Confirm the real Deployment is actually healthy:

   ```
   kubectl get deployment web-web
   kubectl get pods -l app=web-web
   ```

2. Read the chart's real rendered output — not the source templates,
   what Helm actually generated:

   ```
   helm get manifest web
   ```

3. Find the real mismatch between where the Service sends traffic and
   where the container actually listens, fix the chart, and
   `helm upgrade` it.

Get a shell: `docker exec -it breaklab-k8sd-008 bash`
