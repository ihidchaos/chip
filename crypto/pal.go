package crypto

type Spake2pVerifier struct {
	W0 []byte
	ML []byte
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

type P256ECDSASignature []byte

func (s *P256ECDSASignature) Bytes() []byte {
	return nil
}

func SignP256ECDSASignature(plainTex, privateKeyFile []byte) (P256ECDSASignature, error) {
	return P256ECDSASignature{}, nil
}

type P256ECDHDerivedSecret struct {
}

//// AesCcmEncrypt
//// AES 要求的Key长度为16个字节
//// Nonce 13个字节的
//// tagLength = AEADMicLengthBytes
//func AesCcmEncrypt(plainText, key, nonce []byte, tagLength int) (outPut []byte, err error) {
//	if len(key) != SymmetricKeyLengthBytes {
//		return nil, fmt.Errorf("key length mismatch:%d", len(key))
//	}
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return
//	}
//	aesCCM, err := aesccm.NewCCM(block, len(nonce), tagLength)
//	if err != nil {
//		return
//	}
//	outPut = aesCCM.Seal(nil, nonce, plainText, nil)
//	return
//}
//
//func AesCtrEncryptOrDecrypt(inputText, key, iv []byte) (outPut []byte, err error) {
//	//创建一个底层使用AES的密码接口
//	block, err := aes.NewCipher(key)
//	if err != nil {
//		return
//	}
//	//创建一个使用CTR分组接口, iv的长度等于明文分组的长度
//	stream := cipher.NewCTR(block, iv)
//	outPut = make([]byte, len(inputText))
//	stream.XORKeyStream(outPut, inputText)
//	return
//}

type SymmetricKeyContextBase interface {
	KeyHash() uint16
	Release()
	MessageEncrypt(plaintext, nonce, addData []byte, tagSize int) (cipherText []byte, err error)
}
