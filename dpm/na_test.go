package dpm

import (
	"testing"
)

func TestNABase(t *testing.T) {
	assertEqual(t, "127.0.0.1:8080;127.0.0.1:8081", getNAsStr([]NA{
		NA{Host: "127.0.0.1", Port: 8080},
		NA{Host: "127.0.0.1", Port: 8081},
	}), "")
}
