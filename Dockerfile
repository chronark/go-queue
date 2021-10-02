# syntax=docker/dockerfile:1

FROM alpine:3.12 as certs

RUN apk add -U --no-cache ca-certificates=20191127-r4

FROM golang:1.16 as builder

WORKDIR /service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .

ARG SERVICE
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
RUN --mount=type=cache,target=/root/.cache/go-build go build -o svc ./main.go


FROM scratch

COPY --from=certs /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /service/svc /svc

CMD ["/svc"]