package crypto

type Spake2p struct {
}

func (p P256Sha256HkdfHmac) Init(context []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) BeginProver(ws []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) ComputeRoundOne(in, out []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) BeginVerifier(t interface{}, i int, t2 interface{}, i2 int, w0 []byte, ml []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) ComputeRoundTwo(x []byte, out []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) KeyConfirm(verifier []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) GetKeys(out []byte) error {
	return nil
}

func (p P256Sha256HkdfHmac) Clear() {

}

type P256Sha256HkdfHmac struct {
	Spake2p
}

func ComputeWS(iterationCount, setupCode uint32, salt []byte) ([]byte, error) {
	return nil, nil
}
