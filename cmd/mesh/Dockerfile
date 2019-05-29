FROM frolvlad/alpine-glibc

WORKDIR /usr/mesh

COPY mesh-linux-amd64 .
COPY forever.sh .

ENV RPC_PORT=60557           
EXPOSE 60557

RUN chmod +x ./mesh-linux-amd64
RUN chmod +x ./forever.sh

RUN apk add ca-certificates

ENTRYPOINT ./forever.sh