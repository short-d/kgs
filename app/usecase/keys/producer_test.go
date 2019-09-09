package keys

import (
	"errors"
	"fmt"
	"testing"

	"github.com/byliuyang/app/mdtest"
	"github.com/byliuyang/kgs/app/entity"
	"github.com/byliuyang/kgs/app/usecase/keys/gen/gentest"
	"github.com/byliuyang/kgs/app/usecase/repo/repotest"
)

func TestProducer_Produce(t *testing.T) {
	testCases := []struct {
		name      string
		keys      []entity.Key
		expErrors []error
		expKeys   []entity.Key
	}{
		{
			name: "unique keys",
			keys: []entity.Key{
				"ab",
				"bc",
				"ac",
			},
			expErrors: []error{},
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
			expErrors: []error{
				errors.New("key exists: ab"),
				errors.New("key exists: bc"),
			},
			expKeys: []entity.Key{
				"ab",
				"bc",
				"cd",
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			repo := repotest.NewAvailableKeyFake()
			gen := gentest.NewGeneratorStub(testCase.keys)
			logger := mdtest.NewLoggerFake()
			producer := NewProducer(&repo, gen, &logger)
			producer.Produce()

			mdtest.SameElements(t, testCase.expErrors, logger.Errors)
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
	producer := NewProducer(&repo, gen, &logger)
	producer.Produce()

	fmt.Println(logger.Errors)
	fmt.Println(repo.GetKeys())
	// Output:
	// [key exists: ab key exists: bc]
	// [ab bc cd]
}
