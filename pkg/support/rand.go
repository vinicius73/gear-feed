package support

import (
	"crypto/rand"
	"math/big"
)

func RandInt(min, max int32) (int32, error) {
	bMax := big.NewInt(int64(max))

	n, err := rand.Int(rand.Reader, bMax)
	if err != nil {
		return 0, err
	}

	return int32(n.Int64()) + min, nil
}
