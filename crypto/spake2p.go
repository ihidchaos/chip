package crypto

type Spake2p struct {
}

func (p P256Sha256HkdfHmac) Init(context []byte) error {
	return nil
}

type P256Sha256HkdfHmac struct {
	Spake2p
}
