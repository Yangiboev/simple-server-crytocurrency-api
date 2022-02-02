FROM golang:1.17 as builder

RUN mkdir -p $GOPATH/src/github.com/Yangiboev/simple-server-crytocurrency-api
WORKDIR $GOPATH/src/github.com/Yangiboev/simple-server-crytocurrency-api

COPY . ./

RUN cp -r ./config /config 
RUN export CGO_ENABLED=0 && \
    export GOOS=linux && \
    # go mod vendor && \
    make build && \
    mv ./bin/crypto /
FROM alpine
COPY --from=builder config ./config
COPY --from=builder crypto .
ENTRYPOINT [ "/crypto" ]