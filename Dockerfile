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
RUN CGO_ENABLED=0 go build -o /breaklab ./cmd/breaklab

# ---- runtime: needs the real docker CLI, not scratch image ----
# NB: "docker:27-cli-alpine" is not a real tag — the official docker image's
# CLI-only variant is Debian-slim-based ("-cli"), not Alpine, and the current
# major is 29, not 27 (check hub.docker.com/_/docker for the latest).
FROM docker:29-cli AS runtime
WORKDIR /app
COPY --from=frontend-build /src/frontend/dist ./frontend/dist
COPY --from=go-build /breaklab ./breaklab
COPY labs/ ./labs/
COPY templates/web/ ./templates/web/
COPY roadmap.yaml ./roadmap.yaml
EXPOSE 8787
ENTRYPOINT ["./breaklab", "serve", "--host", "0.0.0.0", "--port", "8787"]