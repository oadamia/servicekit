package nanoid

import (
	"strings"
	"testing"
	"unicode"

	"github.com/stretchr/testify/assert"
)

func TestMask(t *testing.T) {
	t.Run("mask", func(t *testing.T) {
		assert := assert.New(t)
		mask := bitmask(len(defaultAlphabet))
		assert.Equal(byte(63), mask)
	})

}

func TestDefault(t *testing.T) {
	prefix := "test"
	g := DefaultGenerator(prefix)

	t.Run("length", func(t *testing.T) {
		assert := assert.New(t)

		id, err := g.Generate()

		assert.NoError(err)
		assert.Equal(defaultLength+len(prefix), len(id))

		assert.Greater(len(defaultAlphabet), 1)
		assert.Less(len(defaultAlphabet), 65)
	})

	t.Run("alphabet unique", func(t *testing.T) {
		assert := assert.New(t)
		unique := make(map[byte]bool)

		for _, v := range defaultAlphabet {
			assert.True(v < unicode.MaxASCII)
			unique[v] = true
		}

		assert.Equal(len(defaultAlphabet), len(unique))
	})

}

func TestCustom(t *testing.T) {
	prefix := "test"

	t.Run("length", func(t *testing.T) {
		assert := assert.New(t)
		g, err := CustomGenerator(prefix, 10, []byte("0123456789"))
		assert.NoError(err)

		id, err := g.Generate()
		assert.NoError(err)
		assert.Equal(10+len(prefix), len(id))

	})

	t.Run("invalid length", func(t *testing.T) {
		assert := assert.New(t)
		_, err := CustomGenerator(prefix, 0, []byte("0123456789"))
		assert.EqualError(err, "length should be greater than 0")

	})

	t.Run("invalid alphabet size", func(t *testing.T) {
		assert := assert.New(t)
		_, err := CustomGenerator(prefix, 10, []byte("0"))
		assert.EqualError(err, "alphabet size should be in a range between 2 and 64")

		_, err = CustomGenerator(prefix, 10, []byte("01234567890123456789012345678901234567890123456789012345678901234567890"))
		assert.EqualError(err, "alphabet size should be in a range between 2 and 64")
	})

	t.Run("alphabet unique and ASCII", func(t *testing.T) {
		assert := assert.New(t)
		_, err := CustomGenerator(prefix, 10, []byte("0123456789abcdeff"))
		assert.EqualError(err, "alphabet should be unique")

	})

	t.Run("alphabet unique and ASCII", func(t *testing.T) {
		assert := assert.New(t)
		_, err := CustomGenerator(prefix, 10, []byte("0123456789abcdef€"))
		assert.EqualError(err, "alphabet should be ASCII")
	})

}

func TestGenerate(t *testing.T) {
	t.Run("generated length", func(t *testing.T) {
		assert := assert.New(t)
		g := DefaultGenerator("")

		id, err := g.Generate()
		assert.NoError(err)
		assert.Equal(26, len(id))
	})

	t.Run("generated characters", func(t *testing.T) {
		assert := assert.New(t)
		g := DefaultGenerator("")

		for i := 0; i < 10; i++ {
			id, err := g.Generate()
			assert.NoError(err)
			for _, char := range id {
				assert.True(strings.ContainsRune(string(defaultAlphabet), char))
			}
		}
	})

	t.Run("generated unique", func(t *testing.T) {
		assert := assert.New(t)
		g := DefaultGenerator("")

		used := make(map[string]bool)
		tries := 100_000

		for i := 0; i < tries; i++ {
			id, err := g.Generate()
			assert.NoError(err)
			assert.False(used[id], "shouldn't return colliding IDs")
			used[id] = true
		}
	})

	t.Run("generated distribution", func(t *testing.T) {
		assert := assert.New(t)
		g := DefaultGenerator("")

		tries := 100_000
		chars := make(map[byte]int)
		for i := 0; i < tries; i++ {
			id, _ := g.Generate()

			for i := 0; i < len(id); i++ {
				chars[id[i]]++
			}
		}

		for _, count := range chars {
			assert.InEpsilon(defaultLength*tries/len(defaultAlphabet), count, .02, "should have flat distribution")
		}
	})
}

func TestNewID(t *testing.T) {
	t.Run("NewID length", func(t *testing.T) {
		assert := assert.New(t)
		prefix := "test"
		g := DefaultGenerator(prefix)

		id, err := g.Generate()
		assert.NoError(err)
		assert.Equal(defaultLength+len(prefix), len(id))
	})

}

func BenchmarkGenerate(b *testing.B) {
	g := DefaultGenerator("")
	for i := 0; i < b.N; i++ {
		_, _ = g.Generate()
	}
}

func BenchmarkNewID(b *testing.B) {
	prefix := "test"
	g := DefaultGenerator(prefix)
	for i := 0; i < b.N; i++ {
		_, _ = g.Generate()
	}
}
