### Rest API

|Field|Type|Mandatory|Description|
|id|number|yes|primary key|
|cmd|string|yes|cmd to be executed|
|output|string|no|cmd output|

```json
{
    "task":{
        "id": 1,
        "cmd": "ls -alh",
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
   -d '{"cmd": "ls -alh"}' \
   .../task

HTTP/1.0 200 OK

{
  "task": {
    "id": 1,
    "cmd": "ls -alh"
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
    "cmd": "ls -alh"
  },
  {
    "id": 1,
    "cmd": "ls -alh"
  },
  {
    "id": 1,
    "cmd": "ls -alh"
  }]
}

```
