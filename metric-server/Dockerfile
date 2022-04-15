FROM golang:1.18 AS builder
LABEL intermediateStageToBeDeleted=true

RUN mkdir -p /build
WORKDIR /build/
COPY . .
# CGO_ENABLED must be disabled to run go binary in Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -mod vendor -o bin/metric-server github.com/manudelca/tp1-distribuidos1/metric-server


FROM busybox:latest
COPY --from=builder /build/bin/metric-server /metric-server
COPY ./metric-server/config.yaml /config.yaml
ENTRYPOINT ["/bin/sh"]