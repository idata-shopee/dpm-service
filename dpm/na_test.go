package dpm

import (
	"github.com/lock-free/obrero"
	"testing"
)

func TestNABase(t *testing.T) {
	assertEqual(t, "127.0.0.1:8080;127.0.0.1:8081", getNAsStr([]obrero.NA{
		obrero.NA{Host: "127.0.0.1", Port: 8080},
		obrero.NA{Host: "127.0.0.1", Port: 8081},
	}), "")
}
