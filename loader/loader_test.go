package loader

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type testContext struct {
	ItemA *int
	ItemB *int
	ItemC *int
}

func TestAdd(t *testing.T) {

	myfuncA := func(ctx context.Context, args interface{}) (interface{}, error) {
		testCtx, ok := args.(testContext)
		if !ok {
			return nil, errors.New("myfuncA invalid args")
		}
		var A = 1234
		testCtx.ItemA = &A
		return testCtx, nil
	}

	Add(myfuncA, 0)
	assert.Equal(t, loads.Len(), 1)

	myfuncB := func(ctx context.Context, args interface{}) (interface{}, error) {
		testCtx, ok := args.(testContext)
		if !ok {
			return nil, errors.New("myfuncB invalid args")
		}
		var B = 1234
		testCtx.ItemB = &B
		return testCtx, nil
	}
	Add(myfuncB, 0)
	assert.Equal(t, loads.Len(), 2)
}

func TestRun(t *testing.T) {

	Flush()

	funcA := func(ctx context.Context, args interface{}) (interface{}, error) {
		testCtx, ok := args.(*testContext)
		if !ok {
			return nil, errors.New("funcA invalid args")
		}
		var A = 1
		testCtx.ItemA = &A
		return testCtx, nil
	}
	Add(funcA, 0)

	funcB := func(ctx context.Context, args interface{}) (interface{}, error) {
		testCtx, ok := args.(*testContext)
		if !ok {
			return nil, errors.New("funcB invalid args")
		}
		var B = *testCtx.ItemA + 1
		testCtx.ItemB = &B
		return testCtx, nil
	}
	Add(funcB, 0)

	funcC := func(ctx context.Context, args interface{}) (interface{}, error) {
		testCtx, ok := args.(*testContext)
		if !ok {
			return nil, errors.New("funcC invalid args")
		}
		var C = *testCtx.ItemB + 1
		testCtx.ItemC = &C
		return testCtx, nil
	}

	Add(funcC, 0)

	testCtx := testContext{
		ItemA: nil,
		ItemB: nil,
		ItemC: nil,
	}

	err := Run(context.Background(), &testCtx)
	assert.Nil(t, err)

	assert.NotNil(t, testCtx.ItemA)
	assert.NotNil(t, testCtx.ItemB)
	assert.NotNil(t, testCtx.ItemC)

	assert.Equal(t, *testCtx.ItemA, 1)
	assert.Equal(t, *testCtx.ItemB, 2)
	assert.Equal(t, *testCtx.ItemC, 3)
}
