package keys

import (
	"fmt"
	"testing"

	"github.com/byliuyang/app/mdtest"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen/gentest"
	"github.com/byliuyang/kgs/app/usecase/repo/repotest"
)

func TestProducer_Produce(t *testing.T) {
	testCases := []struct {
		name    string
		keys    []entity.Key
		hasErr  bool
		expKeys []entity.Key
	}{
		{
			name: "unique keys",
			keys: []entity.Key{
				"ab",
				"bc",
				"ac",
			},
			hasErr: false,
			expKeys: []entity.Key{
				"ab",
				"bc",
				"ac",
			},
		},
		{
			name: "duplicated keys",
			keys: []entity.Key{
				"ab",
				"bc",
				"ab",
				"bc",
				"cd",
			},
			hasErr: true,
			expKeys: []entity.Key{
				"ab",
				"bc",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := repotest.NewAvailableKeyFake()
			gen := gentest.NewGeneratorStub(testCase.keys)
			logger := mdtest.NewLoggerFake()
			producer := NewProducerPersist(repo, gen, &logger)
			err := producer.Produce(uint(len(testCase.expKeys)))

			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
			} else {
				mdtest.Equal(t, nil, err)
			}
			mdtest.SameElements(t, testCase.expKeys, repo.GetKeys())
		})
	}
}

func ExampleProducer_Produce() {
	repo := repotest.NewAvailableKeyFake()
	gen := gentest.NewGeneratorStub(
		[]entity.Key{
			"ab",
			"bc",
			"ab",
			"bc",
			"cd",
		})
	logger := mdtest.NewLoggerFake()
	producer := NewProducerPersist(&repo, gen, &logger)
	err := producer.Produce(1)

	fmt.Println(err)
	fmt.Println(repo.GetKeys())
	// Output:
	// key exists: ab
	// [ab bc]
}
