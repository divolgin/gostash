gostash
=======

logstash client written in Go

## Example
```go
package main

import (
	"github.com/divolgin/gostash"
)

func main() {
	logstash := gostash.NewLogstashClient("logstash_host", "9125", "whale")
	defer logstash.Close()

	metadata := map[string]string{
		"host.name": "myserver",
		"pid":       "345",
	}

	logstash.SendMessage("message logged from myserver by process 345", metadata)
}
```
