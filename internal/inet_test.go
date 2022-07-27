package internal

import (
	"net"
	"testing"
)

func TestIP(t *testing.T) {
	ifs, _ := net.Interfaces()
	for _, i := range ifs {
		addr := GetAddress(i)
		for _, a := range addr {
			t.Log(a.String())
		}
	}
}
