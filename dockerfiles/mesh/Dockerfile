# Note: this must be built from the root of the project with:
#
#     docker build . -f ./cmd/mesh/Dockerfile
#

# mesh-builder produces a statically linked binary
FROM golang:1.14.3-alpine3.11 as mesh-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /0x-mesh

ADD . ./

RUN go build ./cmd/mesh

# Final Image
FROM alpine:3.11

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=mesh-builder /0x-mesh/mesh /usr/mesh/mesh

ENV WS_RPC_ADDR=localhost:60557
EXPOSE 60557

ENV HTTP_RPC_ADDR=localhost:60556
EXPOSE 60556

ENV P2P_TCP_PORT=60558
EXPOSE 60558
ENV P2P_WEBSOCKETS_PORT=60559
EXPOSE 60559

RUN chmod +x ./mesh

ENTRYPOINT ./mesh
