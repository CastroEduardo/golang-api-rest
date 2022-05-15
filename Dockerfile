FROM golang:latest

ENV GOPROXY https://goproxy.cn,direct
WORKDIR $GOPATH/src/github.com/CastroEduardo/golang-api-rest
COPY . $GOPATH/src/github.com/CastroEduardo/golang-api-rest
RUN go build .

EXPOSE 8000
ENTRYPOINT ["./go-gin-example"]
