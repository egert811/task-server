### Rest API

|Field|Type|Mandatory|Description|
|id|number|yes|primary key|
|cmd|string|yes|cmd to be executed|
|args|string[]|no|cmd args|
|output|string|no|cmd output|

```json
{
    "task":{
        "id": 1,
        "cmd": "ls",
        "args": ["-a","-l","-h"],
        "output": "...."
    }
}
```

#### POST /task

Adds a new task to the list

```bash
$ curl -v -i \
   -H "Content-Type: application/json" \
   -X POST \
   -d '{"cmd": "ls", "args":["-a", "-l", "-h"]}' \
   .../task

HTTP/1.0 200 OK

{
  "task": {
    "id": 1,
    "cmd": "ls",
    "args": ["-a","-l","-h"],
  }
}

```

#### GET /task

List tasks registered

```bash
$ curl -v -i \
   -H "Content-Type: application/json" \
   .../task

HTTP/1.0 200 OK

{
  "task":[{
    "id": 1,
    "cmd": "ls",
    "args": ["-a","-l","-h"],
  },
  {
    "id": 2,
    "cmd": "ls",
    "args": ["-a","-l","-h"],
  },
  {
    "id": 3,
    "cmd": "ls",
    "args": ["-a","-l","-h"],
  }]
}

```

#### GET /task/{id}

Get the task details specified

```bash
$ curl -v -i \
   -H "Content-Type: application/json" \
   .../task/1

HTTP/1.0 200 OK

{
  "task":{
    "id": 1,
    "cmd": "ls",
    "args": ["-a","-l","-h"],
    "output": ""
  }
}

```
