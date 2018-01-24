FROM golang:latest

MAINTAINER Allen.Cai "caizz520@gmail.com"

WORKDIR $GOPATH/src/github.com/sundy-li/wechat_brain
ADD . $GOPATH/src/github.com/sundy-li/wechat_brain
RUN go build ./cmd/main.go

EXPOSE 8998

ENTRYPOINT ["./main"]
