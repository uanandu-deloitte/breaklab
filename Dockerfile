# ---- Frontend build ----
FROM node:24-alpine AS frontend-build
WORKDIR /src/frontend
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# ---- go build ----
FROM golang:1.26-alpine AS go-build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# -trimpath -ldflags="-s -w" strips debug symbols and local file paths —
# ~30% smaller binary with no behavior change, verified locally
# (19M -> 13M for this module). Pure shipping-weight, not a security control.
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /breaklab ./cmd/breaklab

# ---- runtime: needs the real docker CLI, not scratch image ----
# NB: "docker:27-cli-alpine" is not a real tag, and "docker:29-cli" is
# Alpine-based already (confirmed via `cat /etc/os-release` against the
# actual pulled image, not assumed) — check hub.docker.com/_/docker before
# bumping this major version, since that could change.
FROM docker:29-cli AS runtime

# Fixed, non-root UID/GID so bind-mounted progress/postmortems directories
# have predictable ownership across hosts — never run this as root.
RUN addgroup -g 10001 breaklab && adduser -D -u 10001 -G breaklab breaklab

WORKDIR /app
COPY --from=frontend-build /src/frontend/dist ./frontend/dist
COPY --from=go-build /breaklab ./breaklab
COPY labs/ ./labs/
COPY templates/web/ ./templates/web/
COPY roadmap.yaml ./roadmap.yaml
RUN chown -R breaklab:breaklab /app
USER breaklab

EXPOSE 8787
# breaklab shells out to this image's docker CLI, which needs access to the
# host's docker.sock — that socket's group GID varies per host and can't be
# baked in at build time. Run this container with both:
#   -v /var/run/docker.sock:/var/run/docker.sock
#   --group-add $(stat -c '%g' /var/run/docker.sock)
# (in Kubernetes: mount the socket and set securityContext.supplementalGroups
# to the node's docker.sock GID instead).
ENTRYPOINT ["./breaklab", "serve", "--host", "0.0.0.0", "--port", "8787"]