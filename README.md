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
--------------------- 
版权声明：本文为CSDN博主「min19900718」的原创文章，遵循CC 4.0 by-sa版权协议，转载请附上原文出处链接及本声明。
原文链接：https://blog.csdn.net/min19900718/article/details/87920254
```
###2. 命令格式
######1. 发布代码json字符串 
```json
    [
        {"name":"xxxx","relative_path":"xxx"},
        {"name":"xxxx","relative_path":"xxx"},
        {"name":"xxxx","relative_path":"xxx"}
    ]

```
######2. 任务切片。publish的通道名称cmd.agentid
```json
    {
      "taskid": 123,
      "executionid": 1,
      "serviceid": "",
      "serviceop": 0,
      "servicename": "",
      "serviceosuser": "",
      "servicemodulename":"",
      "servicedir": "",
      "serviceremotecode": "",
      "servicecodepattern": ["lib","config/static"],
      "servicecustompattern": ["lib/custom.jar","config/template"],
      "servicepidfile": "",
      "servicestartcmd": "",
      "servicestopcmd": ""
      
    }
``` 
######3. 任务切片结果。subscribe的通道名称result.taskid
```json
    {
      "taskid": 123,
      "executionid": 1,
      "rcode": 0,
      "rmsg": "",
      "rsteps": [{"stepnum": 1,"stepmane": "check","stepstate": 0,"stepmsg": "","steptime": ""},
                      {"stepnum": 2,"stepmane": "backup","stepstate": 0,"stepmsg": "","steptime": ""}]  
    }
``` 
######3. agent 重启。publish通道名称grace.agentid
```json
     {
        "agentid":"7422abbe-ada0-46f4-9b60-65c5c2e27a2d",
        "gracecmd": "SIGHUP"
     }
```
