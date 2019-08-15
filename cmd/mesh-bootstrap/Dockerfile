# Note: this must be built from the root of the project with:
#
#     docker build . -f ./cmd/mesh-bootstrap/Dockerfile
#

# mesh-builder produces a statically linked binary
FROM golang:1.12.1-alpine3.9 as mesh-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /go/src/github.com/0xProject/0x-mesh

ADD . ./

RUN make deps-go-no-lockfile

RUN go build ./cmd/mesh-bootstrap

# Final Image
FROM alpine:3.9

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=mesh-builder /go/src/github.com/0xProject/0x-mesh/mesh-bootstrap /usr/mesh/mesh-bootstrap

EXPOSE 60558
EXPOSE 60559

RUN chmod +x ./mesh-bootstrap

ENTRYPOINT ./mesh-bootstrap
