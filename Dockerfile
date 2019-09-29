FROM golang:1.12 AS builder
COPY main.go /go/src/github.com/iMartyn/spotifytwitchsings/
COPY Makefile /go/src/github.com/iMartyn/spotifytwitchsings/
COPY src /go/src/github.com/iMartyn/spotifytwitchsings/src
RUN cd /go/src/github.com/iMartyn/spotifytwitchsings/; make deps
RUN cd /go/src/github.com/iMartyn/spotifytwitchsings/; CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o spotifytwitchsings .
#RUN ls /go/src/github.com/spotifytwitchsings/ -l


FROM scratch
COPY --from=builder /go/src/github.com/iMartyn/spotifytwitchsings/spotifytwitchsings /app/