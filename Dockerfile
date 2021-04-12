FROM golang:1.16.3 AS builder
ADD ./ /src
WORKDIR /src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -a -o stocks-api .

FROM alpine:latest
ENV BADGER_PATH ${BADGER_PATH:-/opt/badger-data}
ENV PORT $PORT
RUN mkdir -p $BADGER_PATH
COPY --from=builder /src/stocks-api /opt/service/stocks-api
WORKDIR /opt/service/
RUN chmod +x /opt/service/stocks-api
EXPOSE $PORT
CMD /opt/service/stocks-api
