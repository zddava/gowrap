package consul

import (
	"context"
	"math/rand"
	"strconv"
	"testing"
)

func TestRegister(t *testing.T) {
	client := NewClient("127.0.0.1", 8500)
	ctx, _ := context.WithCancel(context.Background())

	client.Register(ctx, "test-instance", "test-instance"+strconv.Itoa(rand.Int()), "/health", "127.0.0.1", 8080, nil, nil)

}

func TestUnregister(t *testing.T) {
	client := NewClient("127.0.0.1", 8500)
	ctx, _ := context.WithCancel(context.Background())

	client.Deregister(ctx, "test-instance8719188182406655173")

}
