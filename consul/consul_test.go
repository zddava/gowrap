package consul

import (
	"testing"
)

func TestRegister(t *testing.T) {
	client := NewClient("http://127.0.0.1:8500")
	client.Register("test-instance", "test-instance-1", "/health", "127.0.0.1", 8080, nil, nil)

}

func TestUnregister(t *testing.T) {
	client := NewClient("http://127.0.0.1:8500")
	client.Deregister("test-instance-1")

}
