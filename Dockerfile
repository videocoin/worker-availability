FROM golang:1.18-alpine AS builder
ENV GOPRIVATE=github.com/videocoin/*
ARG BOT_USER="nothing"
ARG BOT_PASSWORD="nothing"
RUN apk add --no-cache ca-certificates git
RUN mkdir /user && \
    echo 'nobody:x:65534:65534:nobody:/:' > /user/passwd && \
    echo 'nobody:x:65534:' > /user/group
RUN git config --global url."https://${BOT_USER}:${BOT_PASSWORD}@github.com".insteadOf "https://github.com"
WORKDIR /src
COPY ./go.mod ./go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o ./build/job ./stats/job

FROM alpine:3.13.6 AS worker-availability
COPY --from=builder /user/group /user/passwd /etc/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/build/job /job
USER nobody:nobody
CMD ["/job"]
