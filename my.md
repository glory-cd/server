###1. 常用命令
######1. protoc编译命令
```
mac:
 protoc --go_out=plugins=grpc:. base.proto
win:
 protoc --plugin=protoc-gen-go=E:\Go\data\bin\protoc-gen-go.exe --go_out=./ rpc.proto
 或
 protoc --go_out=plugins=grpc:. rpc.proto
```
######2. openssl
```
1. openssl genrsa -out server.key 2048
2. openssl req -new -sha256 -key server.key -out server.csr
3. openssl x509 -req -sha256 -in server.csr -signkey server.key -out server.crt -days 3650 

openssl genrsa -out server.key 2048 
openssl req -new -key server.key -subj "/CN=10.30.0.163" -out server.csr
echo subjectAltName = IP:10.30.0.163 > extfile.cnf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -extfile extfile.cnf -out server.crt -days 5000
```