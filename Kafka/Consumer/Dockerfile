FROM golang:1.25-alpine AS builder

RUN apk add build-base

RUN mkdir -p $GOPATH/src/consumer

ADD . $GOPATH/src/consumer
WORKDIR $GOPATH/src/consumer

RUN go mod download
RUN go build -o -tags musl /go/bin/consumer

FROM alpine:latest

COPY broker.txt .

ENV BROKER_HOST=172.17.0.2
ENV BROKER_PORT=9092
ENV BROKER_TOPIC_NAME=test_topic
ENV APP_DEBUG=true
ENV APP_CONTEXT_TIMEOUT=10

COPY --from=builder /go/bin/consumer .
COPY --from=builder /go/src/consumer/entrypoint.sh .
RUN chmod +x ./consumer
RUN chmod +x ./entrypoint.sh


ENTRYPOINT [ "./entrypoint.sh", "./consumer" ]