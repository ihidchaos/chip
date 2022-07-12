package crypto

const (
	KSpake2p_Min_PBKDF_Iterations uint32 = 1000
	KSpake2p_Max_PBKDF_Iterations uint32 = 100000

	kP256_FE_Length                = 32
	kP256_Point_Length             = (2*kP256_FE_Length + 1)
	KSpake2p_Min_PBKDF_Salt_Length = 16
	KSpake2p_Max_PBKDF_Salt_Length = 32

	KSpake2p_VerifierSerialized_Length = kP256_FE_Length + kP256_Point_Length
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
