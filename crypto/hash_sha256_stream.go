package crypto

import (
	"crypto/sha256"
	"hash"
)

type HashSha256Stream struct {
	hash hash.Hash
}

func (hash *HashSha256Stream) Bytes() []byte {
	return hash.hash.Sum(nil)
}

func (hash *HashSha256Stream) AddData(bytes []byte) {
	hash.hash.Write(bytes)
}

func (hash *HashSha256Stream) Begin() {
	hash.hash = sha256.New()
}

func (hash *HashSha256Stream) Finish() []byte {
	return hash.hash.Sum(nil)
}
