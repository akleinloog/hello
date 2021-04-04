# Build Stage
FROM golang:1.16.3-alpine3.13 AS builder
RUN apk add git
WORKDIR /go/src/hello
ADD ./ .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o hello

# Deploy Stage
FROM scratch
ARG VERSION
ARG VCS_REF
LABEL maintainer="arnoud@kleinloog.ch" \
    description="Hello in Go" \
    version=${VERSION} \
    vcs_url="https://github.com/akleinloog/hello" \
    vcs-ref=${VCS_REF}
COPY --from=builder /go/src/hello/hello .
EXPOSE 80 8080
CMD ["./hello", "serve"]