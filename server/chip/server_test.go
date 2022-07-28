package chip

import (
	"log"
	"math/rand"
	"testing"
)

func TestServer_Init(t *testing.T) {
	cid := rand.Uint64()
	nid := rand.Uint64()
	log.Printf("cid: %x  nid: %x", cid, nid)
	//s := sd.makeInstanceName(core.PeerId{}.initCommissionableData(core.CompressedFabricId(cid), core.NodeId(nid)))
	//log.Printf("string: %s", s)
}
