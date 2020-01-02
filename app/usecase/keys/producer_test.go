package keys

import (
	"fmt"
	"github.com/byliuyang/kgs/app/usecase/repo"
	"github.com/byliuyang/kgs/app/usecase/transactional"
	"github.com/byliuyang/kgs/app/usecase/transactional/transactionaltest"
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
			availableKeysRepo := repotest.NewAvailableKeyFake()
			gen := gentest.NewGeneratorStub(testCase.keys)
			logger := mdtest.NewLoggerFake()

			producer := NewProducerPersist(
				func(tx transactional.Transaction) (key repo.AvailableKey, e error) {
					return &availableKeysRepo, nil
				},
				transactionaltest.NewFactoryFake(),
				gen,
				&logger,
			)

			err := producer.Produce(uint(len(testCase.expKeys)))

			if testCase.hasErr {
				mdtest.NotEqual(t, nil, err)
			} else {
				mdtest.Equal(t, nil, err)
			}
			mdtest.SameElements(t, testCase.expKeys, availableKeysRepo.GetKeys())
		})
	}
}

func ExampleProducer_Produce() {
	availableKeysRepo := repotest.NewAvailableKeyFake()
	gen := gentest.NewGeneratorStub(
		[]entity.Key{
			"ab",
			"bc",
			"ab",
			"bc",
			"cd",
		},
	)
	logger := mdtest.NewLoggerFake()

	producer := NewProducerPersist(
		func(tx transactional.Transaction) (key repo.AvailableKey, e error) {
			return &availableKeysRepo, nil
		},
		transactionaltest.NewFactoryFake(),
		gen,
		&logger,
	)

	err := producer.Produce(1)

	fmt.Println(err)
	// Output:
	// key exists: ab
}
