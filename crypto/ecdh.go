package crypto

import (
	"crypto"
	"crypto/ecdh"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"math/big"
	"os"
)

// GenerateEccKey 生成Ecc秘钥对，P256
func GenerateEccKey(privateKeyFile, publicKeyFile string) error {
	//使用ecdsa生成密钥对'
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return err
	}
	//使用 x509对私钥进行序列化
	derText, err := x509.MarshalECPrivateKey(privateKey)
	if err != nil {
		return err
	}
	// 将得到的私钥放入pem.Block结构体中
	derBlock := &pem.Block{
		Type:  "ecdsa private key used elliptic.P256",
		Bytes: derText,
	}
	//使用 pem把私钥编码入文件
	file, err := os.Create(privateKeyFile)
	defer file.Close()
	if err != nil {
		return err
	}
	err = pem.Encode(file, derBlock)
	if err != nil {
		return err
	}
	//从生在的私钥中得到公钥
	publicKey := privateKey.PublicKey
	//使用 x509对公钥进行序列化
	derText, err = x509.MarshalPKIXPublicKey(&publicKey)
	if err != nil {
		return err
	}
	// 将得到的公钥放入pem.Block结构体中
	derBlock = &pem.Block{
		Type:  "ecdsa public key used elliptic.P256",
		Bytes: derText,
	}
	//使用 pem把私钥编码入文件
	_ = file.Close()
	file, err = os.Create(publicKeyFile)
	if err != nil {
		return err
	}
	err = pem.Encode(file, derBlock)
	if err != nil {
		return err
	}
	return nil
}

// EccSign 对数据进行Ecc签名, sha256
func EccSign(plainText []byte, priKey []byte) (sign []byte, err error) {

	//使用x509对Block还原成私钥
	privateKey, err := x509.ParseECPrivateKey(priKey)

	//使用privateKey对Hash值进行签名
	return privateKey.Sign(rand.Reader, sha256.New().Sum(plainText), crypto.SHA256)

}

// EccVerify Ecc签名认证, sha256
func EccVerify(plainText, rText, sTest []byte, publicKeyFile string) (b bool, err error) {

	//从文件中读出公钥 ->[]byte
	file, err := os.Open(publicKeyFile)
	if err != nil {
		return
	}
	info, err := file.Stat()
	if err != nil {
		return
	}
	buf := make([]byte, info.Size())
	_, err = file.Read(buf)
	if err != nil {
		return
	}

	//使用Pem对公钥文件进行解码
	block, _ := pem.Decode(buf)

	//使用x509还原出公钥
	publicInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return
	}
	publicKey := publicInterface.(*ecdsa.PublicKey)

	//对原始数据进行Hash运算
	hash := sha256.Sum256(plainText)

	//签名认证
	var r, s big.Int
	err = r.UnmarshalText(rText)
	err = s.UnmarshalText(sTest)
	b = ecdsa.Verify(publicKey, hash[:], &r, &s)
	return
}

type P256PublicKey ecdh.PublicKey
type P256Keypair ecdh.PrivateKey

func (k *P256Keypair) PubBytes() []byte {
	//p256publicKey := P256PublicKey(k.PublicKey)
	return (*ecdh.PrivateKey)(k).PublicKey().Bytes()
	//return p256publicKey.Marshal()
}

func (k *P256Keypair) PrivateKey() *ecdh.PrivateKey {
	return (*ecdh.PrivateKey)(k)
}

func (k *P256PublicKey) PublicKey() *ecdh.PublicKey {
	return (*ecdh.PublicKey)(k)
}

func (k *P256Keypair) ECDHDeriveSecret(key *P256PublicKey) ([]byte, error) {
	bytes, err := (*ecdh.PrivateKey)(k).ECDH((*ecdh.PublicKey)(key))

	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (k *P256Keypair) ECDSASignMsg(msg []byte) ([]byte, error) {
	sign, err := EccSign(msg, k.PrivateKey().Bytes())
	if err != nil {
		return nil, err
	}
	return sign, nil
}

func GenericP256Keypair() *P256Keypair {
	privateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		panic(err)
	}
	return (*P256Keypair)(privateKey)
}

// UnmarshalPublicKey  接收到的字节序列化成公钥
func UnmarshalPublicKey(data []byte) (*P256PublicKey, error) {
	pubKey, err := ecdh.X25519().NewPublicKey(data)
	if err != nil {
		return nil, err
	}
	return (*P256PublicKey)(pubKey), nil
}
