FROM oven/bun:latest AS builder-bun
WORKDIR /aletis
COPY ./web/package.json ./package.json
COPY ./web/bun.lock ./bun.lock
RUN bun install
COPY ./web .
RUN bun run build

FROM golang:alpine AS modules
WORKDIR /go/src/github.com/AletisSearch/aletis
COPY go.* ./
RUN go mod download
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest \
    && go install github.com/a-h/templ/cmd/templ@latest \
    && go install github.com/tinylib/msgp@latest

FROM modules AS buildergo
WORKDIR /go/src/github.com/AletisSearch/aletis
COPY . .
COPY --from=builder-bun /aletis/dist ./web/dist
RUN sqlc generate \
    && templ generate ./... \
    && go generate ./...
RUN --mount=type=cache,target=/root/.cache/go-build go build -tags 'goexperiment.jsonv2' -ldflags="-s -w" -o ./cmd/aletis/aletis.so ./cmd/aletis

# Docker build
FROM alpine:latest as web

RUN apk --no-cache -U upgrade \
    && apk --no-cache add --upgrade ca-certificates \
    && wget -O /bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 \
    && chmod +x /bin/dumb-init

COPY --from=buildergo /go/src/github.com/AletisSearch/aletis/cmd/aletis/aletis.so /bin/aletis
WORKDIR /etc/aletis/

ENTRYPOINT ["/bin/dumb-init", "--" , "/bin/aletis"]
