# Sample Dockerfile for a fictional Node.js app, "acme-api".
# This file is the subject of the reading exercise — do not "fix" it,
# just analyze it. It is never built; it's just text to read.

# ---------- builder stage ----------
FROM node:20 AS builder
WORKDIR /app

# Mistake #1: the entire app (including source files that change on
# nearly every commit) is copied in before dependencies are installed.
COPY . .

RUN apt-get update && apt-get install -y --no-install-recommends \
    build-essential python3 \
    && rm -rf /var/lib/apt/lists/*
RUN npm ci
RUN npm run build

# ---------- final stage ----------
FROM node:20-slim
WORKDIR /app

# Only these two paths are ever pulled out of the builder stage.
COPY --from=builder /app/dist ./dist
COPY --from=builder /app/node_modules ./node_modules
COPY package.json .

# Block A: a 50MB scratch file created AND removed within the SAME RUN
# instruction.
RUN dd if=/dev/zero of=/tmp/scratch-a.bin bs=1M count=50 2>/dev/null \
    && rm -f /tmp/scratch-a.bin

# Block B: a 50MB scratch file created in one RUN instruction, then
# removed in a separate, LATER RUN instruction.
RUN dd if=/dev/zero of=/tmp/scratch-b.bin bs=1M count=50 2>/dev/null
RUN rm -f /tmp/scratch-b.bin

EXPOSE 3000
CMD ["node", "dist/server.js"]
