package crypto

import "github.com/galenliu/chip/config"

const (
	MaxX509CertificateLength = 600

	P256FELength                 = 32
	P256ECDSASignatureLengthRaw  = 2 * P256FELength
	P256PointLength              = 2*P256FELength + 1
	Sha256HashLength             = 32
	Sha1HashLength               = 20
	SubjectKeyIdentifierLength   = Sha1HashLength
	AuthorityKeyIdentifierLength = Sha1HashLength

	GroupSizeBytes     = P256FELength
	PublicKeySizeBytes = P256PointLength

	AEADMicLengthBytes      = 16
	SymmetricKeyLengthBytes = 16 //对称密钥长度

	ECDHSecretLength     = P256FELength
	ECDSASignatureLength = P256ECDSASignatureLengthRaw

	MaxFELength    = P256FELength
	MaxPointLength = P256PointLength
	MaxHashLength  = Sha256HashLength

	MaxCSRLength = 255

	HashLenBytes = Sha256HashLength

	Spake2pMinPBKDFSaltLength = 16
	Spake2pMaxPBKDFSaltLength = 32
	Spake2pMinPBKDFIterations = 1000
	Spake2pMaxPBKDFIterations = 100000

	P256PrivateKeyLength = GroupSizeBytes
	P256PublicKeyLength  = PublicKeySizeBytes

	AESCCM128KeyLength   = 128 / 8
	AESCCM128BlockLength = AESCCM128KeyLength
	AESCCM128NonceLength = 13
	AESCCM128TagLength   = 16

	MaxSpake2pContextSize     = 1024
	MaxP256keypairContextSize = 512

	EmitDerIntegerWithoutTagOverhead = 1 // 1 sign stuffer
	EmitDerIntegerOverhead           = 3 // NextTag + Length byte + 1 sign stuffer

	MaxHashSha256ContextSize = config.Sha256ContextSize

	Spake2pWSLength                 = P256FELength + 8
	Spake2pVerifierSerializedLength = P256FELength + P256PointLength

	VIDPrefixForCNEncoding  = "Mvid:"
	PIDPrefixForCNEncoding  = "Mpid:"
	VIDAndPIDHexLength      = 2 * 2
	MaxCommonNameAttrLength = 64
)
