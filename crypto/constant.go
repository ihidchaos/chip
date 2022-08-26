package crypto

const (
	AeadMicLengthBytes = 16
	HashLenBytes       = kSha256HashLength

	kSha256HashLength            = 32
	kP256FELength                = 32
	ChipCryptoGroupSizeBytes     = kP256FELength
	kP256PointLength             = 2*kP256FELength + 1
	KSha256HashLength            = 32
	ChipCryptoPublicKeySizeBytes = kP256PointLength

	KP256PublicKeyLength = ChipCryptoPublicKeySizeBytes

	KSpake2pMaxPbkdfIterations uint32 = 100000

	KSpake2pMinPbkdfSaltLength = 16
	KSpake2pMaxPbkdfSaltLength = 32

	KSpake2pVerifierSerializedLength = kP256FELength + kP256PointLength

	ChipCryptoAEADMicLengthBytes uint16 = 16

	SymmetricKeyLengthBytes = 16 //对称密钥长度
)
