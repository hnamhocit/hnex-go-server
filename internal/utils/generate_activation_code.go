package utils

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"
)

func GenerateActivationCode() (string, time.Time) {
	max := big.NewInt(900000)
	expiredAt := time.Now().Add(time.Minute * 1)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "000000", expiredAt
	}

	return fmt.Sprintf("%06d", n.Int64()+100000), expiredAt
}
