FROM mcr.microsoft.com/devcontainers/go:1-1.21

ARG TARGETOS
ARG TARGETARCH

# install nats
ARG NATS_VERSION=2.10.9
RUN curl -L https://github.com/nats-io/nats-server/releases/download/v${NATS_VERSION}/nats-server-v${NATS_VERSION}-${TARGETOS}-${TARGETARCH}.zip -o nats-server.zip
RUN unzip nats-server.zip -d nats-server && cp nats-server/nats-server-v${NATS_VERSION}-${TARGETOS}-${TARGETARCH}/nats-server /usr/bin

# install nats-cli
RUN curl -sf https://binaries.nats.dev/nats-io/natscli/nats@latest | sh && cp nats /usr/bin

