# healer
check your health api,  centralized

# why

There are some places I can not run agents

# run 
```
go run main.go --listen $PORT --config $config
```

$config can be a file on disk or a remote http text file.  Any change on the
configuration content will be reload automatically in the next run

# config

it must be a yaml file which contains health url endpoints, for example
```
- name: firstHealthEndPoint
  healthURL: http://www.exmaple.com/health
- name: hostportcheck
  healthURL: "host:port"
```

the system will try to return a key map with health status injected to the
config item list. The key is determined by those fields in order in the health
item map

- key
- name
- hostname
- ip

# API
 
`Get /`: show all

`Get /:key`: show an item
