#### BUILD STAGE
FROM golang:1.22-alpine AS builder
RUN apk add --update --no-cache build-base
# Set go private
# ENV GOPRIVATE=<go private>
WORKDIR $GOPATH/src/back/

# Copy project files into the builder
COPY . .

RUN go build -a -mod=vendor -ldflags '-linkmode external -w -s -extldflags "-static"' -o ./xm-compagny

#### IMAGE DEFINITION
FROM alpine:3
# Port used by the application, do not override this port
ENV PORT=8080
# Upgrade system dependencies for security patches or bug fixing
RUN apk update && apk upgrade
WORKDIR /app
# Copy the programm with for uid 1000 and gid 1000
COPY --chown=1000:1000 --from=builder /go/src/back/xm-compagny xm-compagny
COPY --chown=1000:1000 --from=builder /go/src/back/config/config.yaml config.yaml

# Expose port, it is just documentation for ops
EXPOSE ${PORT}
# Define a healthcheck to know if container is healthy with curl by getting a response from /status endpoint and parsing result with jq
# HEALTHCHECK --interval=10s --timeout=10s --start-period=20s --retries=3 \ 
# CMD (curl -s http://localhost:${PORT}/status || echo '{"healthy":null}') | jq -e .healthy
# Run with a non-root user
USER 1000:1000

CMD ["./xm-compagny"]