FROM golang:1.12.5
ENV GO111MODULE=on
WORKDIR /go/src/
COPY . ./server/
WORKDIR /go/src/server
RUN go build -mod=vendor -o cdpserver


FROM 192.168.1.77/afis/centos:v1
ENV USER="root"
RUN mkdir -p  /root/cdp
WORKDIR /root/cdp
COPY 75Config/cert_75/server.key ./cert/
COPY 75Config/cert_75/server.crt ./cert/
COPY --from=0 /go/src/server/cdpserver  .
RUN chmod +x cdpserver
CMD  ["./cdpserver"]