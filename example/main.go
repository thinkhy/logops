package main

import (
	"fmt"
	"github.com/thinkhy/logops"
	"time"
)

func main() {
	config := &logops.Config{
		// Address: "http://127.0.0.1:8086",
		Address:  "45.55.21.6:8089",
		Database: "TestDB",
		UseUDP:   true,
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
