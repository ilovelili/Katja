package testutil

import (
	"math/rand"
	"time"

	ci "gx/ipfs/QmP1DfoUjiWH2ZBo1PBH6FupdBucbDepx3HpWmEY6JMUpY/go-libp2p-crypto"
)

func RandTestKeyPair(bits int) (ci.PrivKey, ci.PubKey, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	return ci.GenerateKeyPairWithReader(ci.RSA, bits, r)
}

func SeededTestKeyPair(seed int64) (ci.PrivKey, ci.PubKey, error) {
	r := rand.New(rand.NewSource(seed))
	return ci.GenerateKeyPairWithReader(ci.RSA, 512, r)
}
