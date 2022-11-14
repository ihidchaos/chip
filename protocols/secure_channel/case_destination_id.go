package secure_channel

import (
	"github.com/galenliu/chip/crypto"
	"github.com/galenliu/chip/lib"
)

func GenerateCASEDestinationId(candidateIpkSpan, initiatorRandom []byte, pubKey *crypto.P256PublicKey, fabricId lib.FabricId, nodeId lib.NodeId) (data []byte, err error) {
	return
}
