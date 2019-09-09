package gen

import (
	"errors"

	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/unique"
)

var _ Generator = (*Alphabet)(nil)

// Alphabet represents unique key generator which generates key of length keyLen
// using the characters in the alphabet
type Alphabet struct {
	alphabet []byte
	keyLen   uint
}

// GenerateKeys generate unique keys and send them to keysOut
func (a Alphabet) GenerateKeys(keysOut chan<- entity.Key) {
	if a.keyLen == 0 {
		close(keysOut)
		return
	}

	var chars []byte
	go func() {
		recKey(a.alphabet, chars, a.keyLen, keysOut)
		close(keysOut)
	}()
}

func recKey(alphabet []byte, chars []byte, remaining uint, keys chan<- entity.Key) {
	if remaining == 0 {
		key := entity.Key(chars)
		keys <- key
		return
	}

	for _, char := range alphabet {
		chars = append(chars, char)
		recKey(alphabet, chars, remaining-1, keys)
		chars = chars[:len(chars)-1]
	}
}

// NewAlphabet creates and initialize Alphabet.
// It returns error when  alphabet is either empty or contains duplicated characters
func NewAlphabet(alphabet []byte, keyLen uint) (Alphabet, error) {
	if len(alphabet) < 1 {
		return Alphabet{}, errors.New("alphabet can't be empty")
	}
	if !unique.Characters(alphabet) {
		return Alphabet{}, errors.New("alphabet has to contain unique characters")
	}
	return Alphabet{
		alphabet: alphabet,
		keyLen:   keyLen,
	}, nil
}
