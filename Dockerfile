# mesh-builder produces a statically linked binary
FROM golang:1.12.1-alpine3.9 as bridge-builder


RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers
RUN npm install -g yarn

WORKDIR /go/src/github.com/0xProject/0x-mesh

ADD . ./


RUN make deps

RUN go build ./cmd/demo/sra_bridge

# Final Image
FROM alpine:3.9

RUN apk update && apk add ca-certificates --no-cache

WORKDIR /usr/mesh

COPY --from=bridge-builder /go/src/github.com/0xProject/0x-mesh/sra_bridge /usr/mesh/sra_bridge

RUN chmod +x ./sra_bridge

ENTRYPOINT ./sra_bridge
