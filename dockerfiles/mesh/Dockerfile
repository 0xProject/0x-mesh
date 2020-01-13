# Note: this must be built from the root of the project with:
#
#     docker build . -f ./cmd/mesh/Dockerfile
#

# mesh-builder produces a statically linked binary
FROM golang:1.12.13-alpine3.9 as mesh-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /go/src/github.com/0xProject/0x-mesh

ADD . ./

RUN make gobin

RUN go build ./cmd/mesh

# Final Image
FROM alpine:3.9

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=mesh-builder /go/src/github.com/0xProject/0x-mesh/mesh /usr/mesh/mesh

ENV RPC_ADDR=localhost:60557
EXPOSE 60557

ENV P2P_TCP_PORT=60558
EXPOSE 60558
ENV P2P_WEBSOCKETS_PORT=60559
EXPOSE 60559

RUN chmod +x ./mesh

ENTRYPOINT ./mesh
