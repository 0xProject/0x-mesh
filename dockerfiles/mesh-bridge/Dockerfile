# Note: this must be built from the root of the project with:
#
#     docker build . -f ./cmd/mesh-bridge/Dockerfile
#

# mesh-builder produces a statically linked binary
FROM golang:1.14.3-alpine3.11 as mesh-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /0x-mesh

ADD . ./

RUN go build ./cmd/mesh-bridge

# Final Image
FROM alpine:3.11

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=mesh-builder /0x-mesh/mesh-bridge /usr/mesh/mesh-bridge

RUN chmod +x ./mesh-bridge

ENTRYPOINT ./mesh-bridge
