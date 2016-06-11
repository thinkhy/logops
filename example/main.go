package main

import (
	"github.com/thinkhy/logops"
)

func main() {
	config := &logops.Config{
		// Address:  "http://127.0.0.1:8086",
		Address:  "http://45.55.21.6:8086",
		Database: "TestDB",
	}
	h, err := logops.NewHook(config)
	if err != nil {
		t.Fail()
	}

	module := "workload"
	who := "wellie"
	how := "insert"
	what := "workload"
	h.Write(module, who, how, what)

	module = "testsuite"
	h.Write(module, who, how, what)
	h.tearDown()
	return
}
