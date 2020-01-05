package gen

import (
	"errors"

	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/unique"
)

var _ Generator = (*Alphabet)(nil)

var base62 = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o',
	'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D',
	'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S',
	'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '0', '1', '2', '3', '4', '5', '6', '7',
	'8', '9',
}

// Alphabet represents unique key generator which generates key of length keyLen
// using the characters in the alphabet
type Alphabet struct {
	alphabet []byte
}

// GenerateKeys generate unique keys and send them to keysOut
func (a Alphabet) GenerateKeys(keySize uint, keysOut chan<- entity.Key) {
	if keySize == 0 {
		close(keysOut)
		return
	}

	var chars []byte
	go func() {
		defer close(keysOut)
		recKey(a.alphabet, chars, keySize, keysOut)
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
func NewAlphabet(alphabet []byte) (Alphabet, error) {
	if len(alphabet) < 1 {
		return Alphabet{}, errors.New("alphabet can't be empty")
	}
	if !unique.Characters(alphabet) {
		return Alphabet{}, errors.New("alphabet has to contain unique characters")
	}
	return Alphabet{
		alphabet: alphabet,
	}, nil
}

func NewBase62() []byte {
	return base62
}
