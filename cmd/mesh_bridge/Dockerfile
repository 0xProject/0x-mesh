# Note: this must be built from the root of the project with:
#
#     docker build . -f ./cmd/mesh_bridge/Dockerfile
#

# mesh-builder produces a statically linked binary
FROM golang:1.13.4-alpine3.10 as mesh-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /0x-mesh

ADD . ./

RUN go build ./cmd/mesh_bridge

# Final Image
FROM alpine:3.9

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=mesh-builder /0x-mesh/mesh_bridge /usr/mesh/mesh_bridge

RUN chmod +x ./mesh_bridge

ENTRYPOINT ./mesh_bridge
