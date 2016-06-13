package logops

import (
	"fmt"
	"testing"
	"time"
)

func TestLog(t *testing.T) {
	fmt.Println("TestLog")
	config := &Config{
		// Address:  "http://127.0.0.1:8086",
		Database: "TestDB",
		Address: "http://45.55.21.6:8086",
	}
	h, err := NewHook(config)
	if err != nil {
		t.Fail()
	}

	module := "workload"
	who := "wellie"
	how := "insert"
	what := "workload"
	where := ""
	h.Write(module, who, how, what, where)

	module = "testsuite"
	h.Write(module, who, how, what, where)

	time.Sleep(2 * time.Second)
	h.tearDown()
	return
}
