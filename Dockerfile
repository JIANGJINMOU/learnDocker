FROM golang:1.22-bullseye AS builder
WORKDIR /src
COPY . .
RUN make build

FROM debian:bullseye-slim
RUN apt-get update && apt-get install -y iproute2 iptables && rm -rf /var/lib/apt/lists/*
COPY --from=builder /src/bin/cede /usr/local/bin/cede
ENTRYPOINT ["/usr/local/bin/cede"]
