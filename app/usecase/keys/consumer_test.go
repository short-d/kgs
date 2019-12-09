package keys

import (
	"errors"
	"github.com/byliuyang/app/mdtest"
	"testing"
)

func TestCachedConsumer(t *testing.T) {
	mockConsumer := &fakeConsumer{
		[]fakeConsumerResult{
			{
				keys: []string{"aaaa", "aaab", "aaac", "aaad", "aaae", "aaaf", "aaag", "aaah"},
				err:  nil,
			},
			{
				keys: []string{"baaa", "baab", "baac", "baad", "baae", "baaf", "baag", "baah"},
				err:  nil,
			},
			{
				keys: []string{"caaa", "caab", "caac", "caad", "caae", "caaf", "caag", "caah"},
				err:  nil,
			},
			{
				keys: []string{},
				err:  fakeConsumerError,
			},
		},
	}

	testCases := []struct {
		expected fakeConsumerResult
		count    uint
	}{
		{
			expected: fakeConsumerResult{
				keys: []string{"aaaa", "aaab", "aaac"},
				err:  nil,
			},
			count: 3,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"aaad", "aaae"},
				err:  nil,
			},
			count: 2,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"aaaf"},
				err:  nil,
			},
			count: 1,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"aaag", "aaah", "baaa"},
				err:  nil,
			},
			count: 3,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"baab", "baac", "baad", "baae", "baaf", "baag", "baah", "caaa", "caab", "caac"},
				err:  nil,
			},
			count: 10,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"caad", "caae", "caaf", "caag"},
				err:  nil,
			},
			count: 4,
		},
		{
			expected: fakeConsumerResult{
				keys: []string{"caah"},
				err:  fakeConsumerError,
			},
			count: 10,
		},
	}

	consumer, err := NewCachedConsumer(8, mockConsumer)
	mdtest.Equal(t, err, nil)

	for _, c := range testCases {
		actual, err := consumer.ConsumeInBatch(c.count)

		mdtest.Equal(t, c.expected.keys, actual)
		mdtest.Equal(t, c.expected.err, err)
	}
}

var fakeConsumerError = errors.New("some trouble with delegate")

type fakeConsumerResult struct {
	keys []string
	err  error
}

type fakeConsumer struct {
	set []fakeConsumerResult
}

func (p *fakeConsumer) ConsumeInBatch(maxCount uint) ([]string, error) {
	var res fakeConsumerResult
	n := len(p.set)

	if n == 0 {
		panic("unexpected call")
	}

	res, p.set = p.set[0], p.set[1:]

	return res.keys, res.err
}
