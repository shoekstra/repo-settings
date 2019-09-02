### Step 1 - build app

FROM golang:1.12-alpine AS builder

# Install deps
#   - git required for fetching dependencies
#   - ca-certificates required to call HTTPS endpoints
RUN apk update && \
    apk add --no-cache git ca-certificates && \
    update-ca-certificates

# Create app user
RUN adduser -D -g '' repo-settings

WORKDIR /app
COPY . .

RUN go mod download
RUN go mod verify
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/repo-settings

### Step 2 - build app image

FROM alpine

# Import from builder
COPY --from=builder /bin/repo-settings /bin/repo-settings
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Use an unprivileged user
USER repo-settings

# Run our app
ENTRYPOINT ["/bin/repo-settings"]