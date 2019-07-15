package dpm

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, expect interface{}, actual interface{}, message string) {
	if expect == actual {
		return
	}
	if len(message) == 0 {
		message = fmt.Sprintf("expect %v !=  actual %v", expect, actual)
	}
	t.Fatal(message)
}

func TestWorkerBase(t *testing.T) {
	var workerConf = WorkerConf{
		Workers: []Worker{
			Worker{ServiceType: "a"},
			Worker{ServiceType: "b"},
			Worker{ServiceType: "c"},
		},
	}

	var workers []Worker
	workers = workerConf.GetWorkers("a")

	assertEqual(t, 1, len(workers), "")
	assertEqual(t, "a", workers[0].ServiceType, "")
}
