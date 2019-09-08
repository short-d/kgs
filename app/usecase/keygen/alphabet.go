package keygen

import (
	"errors"

	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/unique"
)

var _ KeyGenerator = (*Alphabet)(nil)

type Alphabet struct {
	alphabet []byte
	keyLen   uint
}

func (a Alphabet) AvailableKeys(keys chan<- entity.Key) {
	if a.keyLen == 0 {
		close(keys)
		return
	}

	var chars []byte
	go func() {
		recKey(a.alphabet, chars, a.keyLen, keys)
		close(keys)
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
