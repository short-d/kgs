package unique

import (
	"fmt"
	"testing"

	"github.com/byliuyang/app/mdtest"
)

func TestCharacters(t *testing.T) {
	testCases := []struct {
		alphabet []byte
		exp      bool
	}{
		{
			alphabet: nil,
			exp:      true,
		},
		{
			alphabet: []byte{},
			exp:      true,
		},
		{
			alphabet: []byte{'a', 'b', 'c', 'd', '1'},
			exp:      true,
		},
		{
			alphabet: []byte{'a', 'b', 'a'},
			exp:      false,
		},
		{
			alphabet: []byte{'a', 'b', 'a', 'b'},
			exp:      false,
		},
		{
			alphabet: []byte{'a', ' ', 'b', ' '},
			exp:      false,
		},
	}

	for _, testCase := range testCases {
		t.Run(string(testCase.alphabet), func(t *testing.T) {
			got := Characters(testCase.alphabet)
			mdtest.Equal(t, testCase.exp, got)
		})
	}
}

func ExampleCharacters() {
	got := Characters([]byte{'a', 'b', 'a', 'b'})
	fmt.Println(got)
	got = Characters([]byte{'a', 'b', 'c'})
	fmt.Println(got)
	// Output:
	// false
	// true
}
