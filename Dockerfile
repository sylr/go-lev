FROM --platform=$BUILDPLATFORM golang:1.23 AS builder

WORKDIR /go/src/go-lev

COPY go*.mod go*.sum ./

RUN go mod download

COPY . .

ARG TARGETOS
ARG TARGETARCH

# Switch shell to bash
SHELL ["bash", "-c"]

RUN git update-index --refresh; GOOS=$TARGETOS GOARCH=$TARGETARCH make build

# ------------------------------------------------------------------------------

FROM scratch

WORKDIR /usr/local/bin

COPY --from=builder /go/src/go-lev/go-lev .

ENTRYPOINT ["/usr/local/bin/go-lev"]
