FROM golang:1.15 as build-env

WORKDIR /go/src/app
ADD . /go/src/app

RUN go get -d -v ./...

RUN go build -o /go/bin/optimizer

FROM gcr.io/distroless/base
COPY --from=build-env /go/bin/optimizer /
ADD config.yaml /

CMD ["/optimizer"]