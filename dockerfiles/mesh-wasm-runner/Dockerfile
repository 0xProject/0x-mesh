# mesh-builder produces a statically linked binary
FROM golang:1.13.4-alpine3.10 as mesh-builder

RUN apk update && apk add ca-certificates nodejs-current npm make git dep gcc build-base musl linux-headers

WORKDIR /0x-mesh

ADD . ./

RUN go build ./cmd/mesh-wasm-runner
WORKDIR /0x-mesh/packages/runner-wasm
RUN npm install
RUN npm run build

# Final Image
FROM chromedp/headless-shell:latest
# Install dumb-init or tini
RUN apt update && apt install dumb-init musl
# or RUN apt install tini
WORKDIR /usr/mesh

COPY --from=mesh-builder /0x-mesh/mesh-wasm-runner /usr/mesh/mesh-wasm-runner
COPY --from=mesh-builder /0x-mesh/packages/runner-wasm/dist /usr/mesh/dist

# Default port for TCP multiaddr
EXPOSE 60558
# Default port for WebSockets multiaddr
EXPOSE 60559

RUN chmod +x ./mesh-wasm-runner

ENTRYPOINT ["dumb-init", "--"]

CMD ["./mesh-wasm-runner"]

