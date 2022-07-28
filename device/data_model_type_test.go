package device

import (
	"fmt"
	"log"
	"math/rand"
	"testing"
)

func TestNodeId(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := rand.Uint64()
		nodeId := NodeId(id)
		t.Log(nodeId.String())
	}
}

func TestInstanceNameRand(t *testing.T) {
	for i := 0; i < 100; i++ {
		id := rand.Uint64()
		//t.Log(strconv.FormatUint(id, 16))
		log.Printf(fmt.Sprintf("%016X", id))
	}
}
