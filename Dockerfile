FROM golang:1.12-alpine3.9 as builder

ADD . $GOPATH/src/github.com/sylr/go-lev
WORKDIR $GOPATH/src/github.com/sylr/go-lev

RUN apk update && apk upgrade && apk add --no-cache alpine-sdk

RUN uname -a && go version
RUN git update-index --refresh; make install

# -----------------------------------------------------------------------------

FROM alpine:3.9

WORKDIR /usr/local/bin
RUN apk --no-cache add ca-certificates
RUN apk update && apk upgrade && apk add --no-cache bash curl
COPY --from=builder "/go/bin/go-lev" .

ENTRYPOINT ["/usr/local/bin/go-lev"]