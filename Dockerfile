FROM oven/bun:latest AS builder-bun
WORKDIR /aletis
COPY ./web .
RUN bun run build

FROM golang:alpine AS buildergo
WORKDIR /go/src/github.com/AletisSearch/aletis
COPY go.* ./
RUN go mod download
COPY . .
COPY --from=builder-bun /aletis/dist ./web/dist
ENV GOEXPERIMENT="jsonv2"
RUN --mount=type=cache,target=/root/.cache/go-build go build -ldflags="-s -w" -o ./cmd/aletis/aletis.so ./cmd/aletis

# Docker build
FROM alpine:latest

RUN apk --no-cache -U upgrade \
    && apk --no-cache add --upgrade ca-certificates \
    && wget -O /bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.5/dumb-init_1.2.5_x86_64 \
    && chmod +x /bin/dumb-init

COPY --from=buildergo /go/src/github.com/AletisSearch/aletis/cmd/aletis/aletis.so /bin/aletis
WORKDIR /etc/aletis/

ENTRYPOINT ["/bin/dumb-init", "--" , "/bin/aletis"]
