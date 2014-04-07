gostash
=======

logstash client written in Go

## Examples

### Calling gostash directly
```go
package main

import (
	"github.com/divolgin/gostash"
)

func main() {
	config := gostash.Config{}
	config.LogstashHost = "logstash_host"
	config.LogstashPort = "9125"
	config.InputType = "mylogformat"

	logstash := gostash.NewLogstashClient(&config)
	defer logstash.Close()

	metadata := map[string]string{
		"host.name": "myserver",
		"pid":       "345",
	}

	logstash.SendMessage("message logged from myserver by process 345", metadata)
}
```

### Logging with loggo
```go
package main

import (
	"github.com/divolgin/gostash"
	"github.com/lookio/loggo"
)

func main() {
	config := gostash.Config{}
	config.LogstashHost = "logstash_host"
	config.LogstashPort = "9125"
	config.InputType = "mylogformat"

	logstash := gostash.NewLogstashClient(&config)
	loggoWriter := loggo.NewSimpleWriter(logstash, logstash.Formatter())
	loggo.RegisterWriter("logstash", loggoWriter, loggo.DEBUG)

	logger := loggo.GetLogger("")
	logger.Debugf("this message will be sent to logstash")
	logger.Errorf("this message will be printed to console and sent to logstash")
}
```
