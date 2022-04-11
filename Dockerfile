FROM golang:1.18-bullseye as builder
RUN apt-get update && \
    apt-get -y install git openssh-client make curl && \
    rm -rf /var/lib/apt/lists/*

# COPY only the go mod files for efficient caching
COPY go.mod go.sum /go/src/github.com/akumor/keepassnotifier/
WORKDIR /go/src/github.com/akumor/keepassnotifier

# Pull dependencies
RUN go mod download

# COPY the rest of the source code
COPY . /go/src/github.com/akumor/keepassnotifier/

# This 'linux_compile' target should compile binaries to the /artifacts directory
# The main entrypoint should be compiled to /artifacts/keepassnotifier
RUN make linux_compile

# update the PATH to include the /artifacts directory
ENV PATH="/artifacts:${PATH}"

FROM ubuntu:18.04
LABEL org.opencontainers.image.source https://github.com/akumor/keepassnotifier

COPY --from=builder /artifacts /bin

# Ensure the latest CA certs are present to authenticate SSL connections.
RUN apt-get update && \
    apt-get -y install ca-certificates && \
    rm -rf /var/lib/apt/lists/*

RUN groupadd -r keepassnotifier && \
    useradd -ms /bin/bash keepassnotifier -g keepassnotifier
USER keepassnotifier

CMD ["keepassnotifier"]
