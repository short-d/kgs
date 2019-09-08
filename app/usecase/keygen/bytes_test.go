package keygen

import (
	"fmt"
	"testing"

	"github.com/byliuyang/app/mdtest"
	"github.com/byliuyang/kgs/app/entity"
)

func TestNewCharacters(t *testing.T) {
	testCases := []struct {
		alphabet  []byte
		expHasErr bool
	}{
		{
			alphabet:  nil,
			expHasErr: true,
		},
		{
			alphabet:  []byte{},
			expHasErr: true,
		},
		{
			alphabet:  []byte{'a', 'b', 'c', 'd', '1'},
			expHasErr: false,
		},
		{
			alphabet:  []byte{'a', 'b', 'a'},
			expHasErr: true,
		},
		{
			alphabet:  []byte{'a', 'b', 'a', 'b'},
			expHasErr: true,
		},
		{
			alphabet:  []byte{'a', ' ', 'b', ' '},
			expHasErr: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.alphabet), func(t *testing.T) {
			_, err := NewCharacters(testCase.alphabet, 5)
			if testCase.expHasErr {
				mdtest.NotEqual(t, nil, err)
				return
			}
			mdtest.Equal(t, nil, err)
		})
	}
}

func TestCharacters_AvailableKeys(t *testing.T) {
	testCases := []struct {
		alphabet []byte
		keyLen   uint
		expKeys  []entity.Key
	}{
		{
			alphabet: []byte{'a'},
			keyLen:   0,
			expKeys:  []entity.Key{},
		},
		{
			alphabet: []byte{'a'},
			keyLen:   1,
			expKeys:  []entity.Key{"a"},
		},
		{
			alphabet: []byte{'a'},
			keyLen:   2,
			expKeys:  []entity.Key{"aa"},
		},
		{
			alphabet: []byte{'a'},
			keyLen:   3,
			expKeys:  []entity.Key{"aaa"},
		},
		{
			alphabet: []byte{'a', 'b'},
			keyLen:   0,
			expKeys:  []entity.Key{},
		},
		{
			alphabet: []byte{'a', 'b'},
			keyLen:   1,
			expKeys:  []entity.Key{"a", "b"},
		},
		{
			alphabet: []byte{'a', 'b'},
			keyLen:   2,
			expKeys:  []entity.Key{"aa", "ab", "ba", "bb"},
		},
		{
			alphabet: []byte{'a', 'b'},
			keyLen:   3,
			expKeys: []entity.Key{
				"aaa",
				"aab",
				"aba",
				"abb",
				"baa",
				"bab",
				"bba",
				"bbb",
			},
		},
	}

	for _, testCase := range testCases {
		name := fmt.Sprintf("%s %d", testCase.alphabet, testCase.keyLen)
		t.Run(name, func(t *testing.T) {
			chars, err := NewCharacters(testCase.alphabet, testCase.keyLen)
			mdtest.Equal(t, nil, err)

			availableKeys := make(chan entity.Key)
			chars.AvailableKeys(availableKeys)

			var gotKeys = collectKeys(availableKeys)
			mdtest.SameElements(t, testCase.expKeys, gotKeys)
		})
	}
}

func ExampleCharacters_AvailableKeys() {
	chars, _ := NewCharacters([]byte{'a', 'b'}, 3)
	keyChan := make(chan entity.Key)
	chars.AvailableKeys(keyChan)
	keys := collectKeys(keyChan)
	fmt.Println(keys)
	// Output: [aaa aab aba abb baa bab bba bbb]
}

func collectKeys(keyChan <-chan entity.Key) []entity.Key {
	var keys []entity.Key
	for key := range keyChan {
		keys = append(keys, key)
	}
	return keys
}
