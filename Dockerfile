FROM node:25 AS ui_builder

WORKDIR /build

COPY ui/package.json        .
COPY ui/package-lock.json   .

RUN npm install --no-audit

COPY ui .

RUN npm run build

FROM golang:1-alpine3.23 AS builder

RUN apk update && apk add build-base

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:3

WORKDIR /app

COPY --from=ui_builder /build/dist ui/dist
COPY --from=builder /build/app app

EXPOSE 8080

CMD [ "/app/app" ]

### build ###
# export docker_http_proxy=http://host.docker.internal:1080
# docker build --platform linux/amd64 --build-arg http_proxy=$docker_http_proxy --build-arg https_proxy=$docker_http_proxy -f Dockerfile -t allape/golang:latest .
