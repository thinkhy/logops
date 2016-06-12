# logops

[![Go Report
Card](https://goreportcard.com/badge/github.com/thinkhy/logops)](https://goreportcard.com/report/github.com/thinkhy/logops)

Golang package for sending operations log to InfluxDB


Installation
---------------
  
```shell
go get github.com/thinkhy/logops
```

Example
----------

```golang
package main

import (
	"fmt"
	"github.com/thinkhy/logops"
	"time"
)

func main() {
	config := &logops.Config{
		Address: "http://45.55.21.6:8086",
		// Address:  "45.55.21.6:8089",
		Database: "TestDB",
		// UseUDP:   true,
	}
	h, err := logops.NewHook(config)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	module := "workload"
	who := "wellie"
	how := "insert"
	what := "workload"
	h.Write(module, who, how, what)

	module = "testsuite"
	h.Write(module, who, how, what)
	time.Sleep(3 * time.Second)
	return
}
```

Activities (Timeline)
-----------------------
  - When: datatime
  - Who:  operator
  - How:  action
  - What: entry

Concepts
------------------

### Concurrency

We use a non-blocking model for the hook. Each time an entry is fired, we create an InfluxDB point and put that through a buffered channel. It is then picked up by a worker goroutine that will handle the current batch and write everything once the batch hit the size or time limit. 

### Configuration

We chose to take in as a parameter the InfluxDB client so that we have greater flexibility over the way we connect to the Influx Server.
All the defaults and configuration options can be found in the GoDoc.

### Database Handling

The database should have been previously created. This is by design, as we want to avoid people generating lots of databases.
When passing an empty string for the InfluxDB database name, we default to "logrus" as the database name.

### Message Field

We will insert your message into InfluxDB with the field message.
