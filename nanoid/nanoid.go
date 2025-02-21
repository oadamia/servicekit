package nanoid

import (
	"crypto/rand"
	"fmt"
	"math/bits"
	"unicode"
)

// Note: The defaultAlphabet size should be in a range between 2 and 64
var defaultAlphabet = []byte("0123456789abcdefghijklmnopqrstuvwxyz")

const (
	// defaultLength is the defaultLength of the generated random ID
	defaultLength = 26
)

type Generator interface {
	Generate() (string, error)
}

type NanoGenerator struct {
	length   int
	alphabet []byte
	prefix   string
	step     int
	bitmask  byte
}

func DefaultGenerator(prefix string) *NanoGenerator {
	gen, _ := CustomGenerator(prefix, defaultLength, defaultAlphabet)

	return gen
}

func CustomGenerator(prefix string, length int, alphabet []byte) (*NanoGenerator, error) {
	err := validateInput(length, alphabet)
	if err != nil {
		return nil, err
	}

	return &NanoGenerator{
		prefix:   prefix,
		length:   length,
		alphabet: alphabet,
		step:     (length / 5) * 8,
		bitmask:  bitmask(len(alphabet)),
	}, nil
}

func validateInput(length int, alphabet []byte) error {
	if length < 1 {
		return fmt.Errorf("length should be greater than 0")
	}

	if len(alphabet) < 2 || len(alphabet) > 64 {
		return fmt.Errorf("alphabet size should be in a range between 2 and 64")
	}

	unique := make(map[byte]bool)
	for _, v := range alphabet {
		if v >= unicode.MaxASCII {
			return fmt.Errorf("alphabet should be ASCII")
		}
		unique[v] = true
	}

	if len(alphabet) != len(unique) {
		return fmt.Errorf("alphabet should be unique")
	}
	return nil
}

// bitmask used to obtain bits from the random bytes
// getBitmask generates bit mask used to obtain bits from the random bytes that are used to get index of random character
// from the alphabet. Example: if the alphabet has 6 = (110)_2 characters it is sufficient to use mask 7 = (111)_2
func bitmask(length int) byte {
	x := uint8(length) - 1
	mask := 1<<uint(8-bits.LeadingZeros8(x)) - 1
	return byte(mask)
}

func (g *NanoGenerator) Generate() (string, error) {
	id, err := g.generate()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s%s", g.prefix, id), nil
}

func (g *NanoGenerator) generate() (string, error) {
	id := make([]byte, g.length)
	bytes := make([]byte, g.step)

	for j := 0; ; {
		_, err := rand.Read(bytes)
		if err != nil {
			return "", err
		}

		for i := 0; i < g.step; i++ {
			idx := bytes[i] & g.bitmask
			if idx < byte(len(g.alphabet)) {
				id[j] = g.alphabet[idx]
				j++
				if j == g.length {
					return string(id), nil
				}
			}
		}
	}
}
