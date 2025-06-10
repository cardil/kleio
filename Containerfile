FROM docker.io/library/golang:1.24 as builder

COPY . /code

WORKDIR /code
RUN go build -o /bin/collector ./cmd/collector

FROM registry.access.redhat.com/ubi9/ubi-minimal

COPY --from=builder /bin/collector /usr/sbin/collector

EXPOSE 514

ENTRYPOINT ["/usr/sbin/collector"]
