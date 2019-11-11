FROM golang:1.12 AS builder
LABEL maintainer="Daniel Lynch <danplynch@gmail.com>"
RUN mkdir -p /go/src/github.com/randomtask1155/hqserver
WORKDIR $GOPATH/src/github.com/randomtask1155/hqserver
ENV GOPATH=/go
ENV PATH=$GOPATH/bin:$PATH
ENV PORT=8080
ADD . .

RUN GOOS=linux GOARCH=arm GOARM=7 go build -o hqserver .

FROM scratch
COPY --from=builder /go/src/github.com/randomtask1155/hqserver/hqserver /go/bin/hqserver
EXPOSE 8080
ENTRYPOINT ["/go/bin/hqserver"]
