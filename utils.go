package bw6

import (
	"math/big"
)

func bigFromHex(hex string) *big.Int {
	if len(hex) > 1 && hex[:2] == "0x" {
		hex = hex[2:]
	}
	n, _ := new(big.Int).SetString(hex, 16)
	return n
}

func computeNaf(k *big.Int) []int8 {
	// Algorithm 3.30 Computing the NAF of a positive integer p.98
	// Guide to Elliptic Curve Cryptography Menezes
	result := []int8{}
	i := 0
	for k.Cmp(big.NewInt(1)) != -1 {
		var k_i *big.Int
		if k.Bit(0) != 0 {
			tmp := new(big.Int).Mod(k, big.NewInt(4))
			k_i = new(big.Int).Sub(big.NewInt(2), tmp)
			k = new(big.Int).Sub(k, k_i)
		} else {
			k_i = big.NewInt(0)
		}
		// k to int8
		tmp := int8(k_i.Int64())
		result = append(result, tmp)
		k = new(big.Int).Div(k, big.NewInt(2))
		i += 1
	}
	return result
}
