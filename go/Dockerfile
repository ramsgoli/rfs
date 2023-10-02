ARG  APP


FROM golang:1.20-buster as builder
ARG  APP

# Create and change to the app directory.
WORKDIR /opt

# Retrieve application dependencies.
# This allows the container build to reuse cached dependencies.
# Expecting to copy go.mod and if present go.sum.
COPY go.* ./
RUN go mod download

# Copy local code to the container image.
COPY . ./

# Build the binary.
RUN go build -o ./bin/${APP} ./cmd/${APP}

# Use the official Debian slim image for a lean production container.
# https://hub.docker.com/_/debian
# https://docs.docker.com/develop/develop-images/multistage-build/#use-multi-stage-builds
FROM debian:buster-slim
ARG APP
ENV APP ${APP}
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /opt/bin/${APP} /opt/bin/${APP}

# allow the binary to run as PID 1 in order to intercept signals
RUN echo "#!/bin/bash \n exec /opt/bin/${APP}" > ./entrypoint.sh
RUN chmod +x ./entrypoint.sh
ENTRYPOINT ["./entrypoint.sh"]

# [END run_helloworld_dockerfile]
# [END cloudrun_helloworld_dockerfile]