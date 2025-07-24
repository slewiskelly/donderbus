ARG GOVERSION=1

FROM golang:${GOVERSION}-alpine AS builder

ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/go/pkg/mod

WORKDIR /src

RUN apk --no-cache add git

COPY ./go.* ./
RUN --mount=type=cache,target=${GOMODCACHE} go mod download

COPY . ./
RUN --mount=type=cache,target=${GOCACHE} \
    --mount=type=cache,target=${GOMODCACHE} \
    CGO_ENABLED=0 go build -o /bin/ ./cmd/*

FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=builder /bin/* /bin/

ENTRYPOINT ["donderbus"]
