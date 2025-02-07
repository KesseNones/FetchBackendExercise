from alpine:latest

from golang:1.23 as builder

workdir /root

copy main.go ./

copy go.mod ./

run go build -o server-app .

entrypoint ["/root/server-app"]
