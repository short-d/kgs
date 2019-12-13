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

	testCases := []struct {
		name       string
		expected   fakeConsumerResult
		skipBefore uint
		count      uint
	}{
		{
			name: "should load keys and return three items from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaaa", "aaab", "aaac"},
				err:  nil,
			},
			skipBefore: 0,
			count: 3,
		},
		{
			name: "should load keys twice and return the entire buffer #2",
			expected: fakeConsumerResult{
				keys: []string{"baaa", "baab", "baac", "baad", "baae", "baaf", "baag", "baah"},
				err:  nil,
			},
			skipBefore: 8,
			count: 8,
		},
		{
			name: "should return two items from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaad", "aaae"},
				err:  nil,
			},
			skipBefore: 3,
			count: 2,
		},
		{
			name: "should return one item from the buffer #1",
			expected: fakeConsumerResult{
				keys: []string{"aaaf"},
				err:  nil,
			},
			skipBefore: 5,
			count: 1,
		},
		{
			name: "should return the first two items from the buffer #1, load new keys and return one key from the buffer #2",
			expected: fakeConsumerResult{
				keys: []string{"aaag", "aaah", "baaa"},
				err:  nil,
			},
			skipBefore: 6,
			count: 3,
		},
		{
			name: "should return the first seven items from the buffer #2, load new keys and return three items from the buffer #3",
			expected: fakeConsumerResult{
				keys: []string{"baab", "baac", "baad", "baae", "baaf", "baag", "baah", "caaa", "caab", "caac"},
				err:  nil,
			},
			skipBefore: 9,
			count: 10,
		},
		{
			name: "should return four items from the buffer #3",
			expected: fakeConsumerResult{
				keys: []string{"caad", "caae", "caaf", "caag"},
				err:  nil,
			},
			skipBefore: 19,
			count: 4,
		},
		{
			name: "should return the last item from the buffer #3 and return an error",
			expected: fakeConsumerResult{
				keys: []string{"caah"},
				err:  fakeConsumerError,
			},
			skipBefore: 23,
			count: 10,
		},
		{
			name: "should return empty list because there are no more keys",
			expected: fakeConsumerResult{
				keys: []string{},
				err:  nil,
			},
			skipBefore: 24,
			count: 10,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
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

			consumer, err := NewCachedConsumer(8, mockConsumer)
			mdtest.Equal(t, err, nil)

			_, err = consumer.ConsumeInBatch(testCase.skipBefore)
			mdtest.Equal(t, err, nil)

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
