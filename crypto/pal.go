package crypto

const (
	KSpake2pMaxPbkdfIterations uint32 = 100000
	kp256FeLength                     = 32
	kP256PointLength                  = 2*kp256FeLength + 1
	KSpake2pMinPbkdfSaltLength        = 16
	KSpake2pMaxPbkdfSaltLength        = 32

	KSpake2pVerifierSerializedLength = kp256FeLength + kP256PointLength

	ChipCryptoAeadMicLengthBytes uint16 = 16
)

type Spake2pVerifier struct {
}

func (v Spake2pVerifier) Deserialize(verifier []byte) error {
	return nil
}

func (v Spake2pVerifier) Generate(count uint32, span []byte, passcode uint32) error {
	return nil
}

func (v Spake2pVerifier) Serialize() ([]byte, error) {
	return nil, nil
}

type P256ECDSASignature struct {
}

func SignP256ECDSASignature(plainTex, privateKeyFile []byte) (P256ECDSASignature, error) {
	return P256ECDSASignature{}, nil
}

type P256ECDHDerivedSecret struct {
}
