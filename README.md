#1 Summary
   cdp server and agent work together to complete deployment, upgrade, etc
   cdp server is similar to a center, which commands the agent to complete the specified task

#2 require
* redis
* etcd
* certificate

#3 message format
##2.1 task work. 
server can publish task message to agent, channel name is "cmd.node-id", node-id is according to the actual situation.

```json
[{      "taskid": 123,
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
        "servicestopcmd": ""}]
```
##2.2 task work result
server receive work result, so server need subscribe one channel when it publish task work to agent.
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
##2.3 special task
restart agent is one special task, just publish message in channel "grace.node-id"
```json
  {
        "agentid":"7422abbe-ada0-46f4-9b60-65c5c2e27a2d",
        "gracecmd": "SIGHUP"
     }
```