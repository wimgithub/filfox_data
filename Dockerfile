FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/filfox_data
COPY . $GOPATH/src/filfox_data
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./filfox_data"]
