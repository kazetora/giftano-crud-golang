FROM golang:1.10-alpine AS builder
ARG SVC_NAME

RUN mkdir /go/src/giftano-crud-golang
ADD . /go/src/giftano-crud-golang/
WORKDIR /go/src/giftano-crud-golang
RUN apk update \
    cd /go/src/giftano-crud-golang \
    && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"'  -o build/$SVC_NAME cmd/$SVC_NAME.go \
    && mv build/$SVC_NAME /exe

FROM scratch
COPY --from=builder /exe /
ENTRYPOINT ["/exe"]
