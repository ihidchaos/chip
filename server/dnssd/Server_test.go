package dnssd

import (
	"crypto/rand"
	"fmt"
	"io"
	"testing"
)

func TestDnssd(t *testing.T) {
	bytes8 := make([]byte, 8)
	_, _ = io.ReadFull(rand.Reader, bytes8)
	t.Log(fmt.Sprintf("%016X", bytes8))

}
