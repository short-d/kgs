package keys

import (
	"fmt"
	"testing"

	"github.com/short-d/app/mdtest"
	"github.com/short-d/kgs/app/entity"
	"github.com/short-d/kgs/app/usecase/keys/gen/gentest"
	"github.com/short-d/kgs/app/usecase/repo"
	"github.com/short-d/kgs/app/usecase/repo/repotest"
	"github.com/short-d/kgs/app/usecase/transactional"
	"github.com/short-d/kgs/app/usecase/transactional/transactionaltest"
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
			logger := mdtest.NewLoggerFake(mdtest.FakeLoggerArgs{})

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
	logger := mdtest.NewLoggerFake(mdtest.FakeLoggerArgs{})

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
