FROM node:25 AS builder_ui_common

WORKDIR /build/common

COPY ui/common .

RUN npm install --no-audit

FROM node:25 AS builder_ui_admin

WORKDIR /build/common

COPY --from=builder_ui_common /build/common .

WORKDIR /build/ui

COPY ui/admin/package.json        .
COPY ui/admin/package-lock.json   .

RUN npm install --no-audit

COPY ui/admin .

RUN npm run build

FROM node:25 AS builder_ui_client

WORKDIR /build/common

COPY --from=builder_ui_common /build/common .

WORKDIR /build/ui

COPY ui/client/package.json        .
COPY ui/client/package-lock.json   .

RUN npm install --no-audit

COPY ui/client .

RUN npm run build

FROM golang:1-alpine3.23 AS builder

RUN apk add --no-cache build-base

WORKDIR /build

COPY go.mod go.mod
COPY go.sum go.sum
RUN go mod download

COPY . .
RUN go build -o app .

FROM alpine:3

RUN apk add --no-cache ffmpeg exiftool

WORKDIR /app

COPY --from=builder_ui_admin /build/ui/dist ui/admin/dist
COPY --from=builder_ui_client /build/ui/dist ui/client/dist
COPY --from=builder /build/app app

EXPOSE 8080
EXPOSE 8888

CMD [ "/app/app" ]

### BUILD ###
# export x_docker_http_proxy="http://host.docker.internal:1080"
# export x_docker_image_name="allape/projectname"
# export x_docker_registry_prefix="docker-registry.lan.allape.cc/"
# export x_docker_registry_image_name="$x_docker_registry_prefix$x_docker_image_name"

# docker build --platform linux/arm64 --build-arg http_proxy=$x_docker_http_proxy --build-arg https_proxy=$x_docker_http_proxy -f Dockerfile -t "$x_docker_image_name-arm64" .
# docker build --platform linux/amd64 --build-arg http_proxy=$x_docker_http_proxy --build-arg https_proxy=$x_docker_http_proxy -f Dockerfile -t $x_docker_image_name .
# docker tag $x_docker_image_name $x_docker_registry_image_name && docker push $x_docker_registry_image_name

# sudo docker pull $x_docker_registry_image_name && sudo docker tag $x_docker_registry_image_name $x_docker_image_name
# sudo docker compose -f docker.compose.yaml up -d
