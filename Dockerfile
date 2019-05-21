ARG GOLANG_VERSION=latest
FROM golang:${GOLANG_VERSION} AS builder

ADD . /tls-echo
WORKDIR /tls-echo

ENV CGO_ENABLED=0
RUN go build


FROM alpine:latest
COPY --from=builder /tls-echo/tls-echo /bin/tls-echo

RUN mkdir /etc/tls-echo
VOLUME ["/etc/tls-echo"]

EXPOSE 443
WORKDIR /etc/tls-echo
ENTRYPOINT ["/bin/tls-echo"]
