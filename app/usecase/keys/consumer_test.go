package keys

import (
	"errors"
	"testing"

	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/repo/repotest"

	"github.com/byliuyang/app/mdtest"
)

func TestCachedConsumer(t *testing.T) {
	keys := []entity.Key{
		"aaaa", "aaab", "aaac", "aaad", "aaae", "aaaf", "aaag", "aaah",
		"baaa", "baab", "baac", "baad", "baae", "baaf", "baag", "baah",
		"caaa", "caab", "caac", "caad", "caae", "caaf", "caag", "caah",
	}

	availableKeysRepo := repotest.NewAvailableKeyFake()
	allocatedKeysRepo := repotest.NewAllocatedKeyFake()

	for _, key := range keys {
		err := availableKeysRepo.Create(key)
		mdtest.Equal(t, nil, err)
	}

	mockConsumer := NewConsumerPersist(
		&availableKeysRepo,
		&allocatedKeysRepo,
	)

	testCases := []struct {
		name     string
		expected fakeConsumerResult
		count    uint
	}{
		{
			name: "Should load keys and return three items from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaaa", "aaab", "aaac"},
				err:  nil,
			},
			count: 3,
		},
		{
			name: "Should return two items from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaad", "aaae"},
				err:  nil,
			},
			count: 2,
		},
		{
			name: "Should return one item from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaaf"},
				err:  nil,
			},
			count: 1,
		},
		{
			name: "Should return the first two items from the buffer #1, load new keys and return one key from the buffer #2",
			expected: fakeConsumerResult{
				keys: []string{"aaag", "aaah", "baaa"},
				err:  nil,
			},
			count: 3,
		},
		{
			name: "Should return the first seven items from the buffer #2, load new keys and return three items from the buffer #3",
			expected: fakeConsumerResult{
				keys: []string{"baab", "baac", "baad", "baae", "baaf", "baag", "baah", "caaa", "caab", "caac"},
				err:  nil,
			},
			count: 10,
		},
		{
			name: "Should return four items from the buffer #3",
			expected: fakeConsumerResult{
				keys: []string{"caad", "caae", "caaf", "caag"},
				err:  nil,
			},
			count: 4,
		},
		{
			name: "Should return the last item from the buffer #3 and return an error",
			expected: fakeConsumerResult{
				keys: []string{"caah"},
				err:  fakeConsumerError,
			},
			count: 10,
		},
	}

	consumer, err := NewCachedConsumer(8, mockConsumer)
	mdtest.Equal(t, err, nil)

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			allocatedKeysRepo.FakeError(testCase.expected.err)
			actual, err := consumer.ConsumeInBatch(testCase.count)

			mdtest.Equal(t, testCase.expected.keys, actual)
			mdtest.Equal(t, testCase.expected.err, err)
		})
	}
}

var fakeConsumerError = errors.New("some trouble with delegate")

type fakeConsumerResult struct {
	keys []string
	err  error
}
