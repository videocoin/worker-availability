FROM golang:1.14-alpine as builder

RUN apk add --no-cache make gcc musl-dev linux-headers git
WORKDIR /build/
ADD . .
COPY cloud-api@v1.1.3 /tmp/cloud-api

RUN make build

FROM alpine:latest

COPY --from=builder /build/build/job /

WORKDIR /
CMD ["/job"]