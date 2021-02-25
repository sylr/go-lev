FROM golang:1.16 as builder

ADD . $GOPATH/src/github.com/sylr/go-lev
WORKDIR $GOPATH/src/github.com/sylr/go-lev

RUN uname -a && go version
RUN git update-index --refresh; make install

# -----------------------------------------------------------------------------

FROM scratch

WORKDIR /usr/local/bin

COPY --from=builder "/go/bin/go-lev" .

ENTRYPOINT ["/usr/local/bin/go-lev"]