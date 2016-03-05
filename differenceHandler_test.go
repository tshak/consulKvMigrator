package consulKvMigrator

import "testing"
import "github.com/stretchr/testify/assert"

func TestChainHandlers(t *testing.T) {
	handlerA := &testableHandler{}
	handlerB := &testableHandler{}
	handlerC := &testableHandler{}
	head := chainHandlers([]differenceHandler{handlerA, handlerB, handlerC})

	assert.Equal(t, head, handlerA)
	assert.Equal(t, handlerA.nextHandler, handlerB)
	assert.Equal(t, handlerB.nextHandler, handlerC)
	assert.Nil(t, handlerC.nextHandler)
}

type testableHandler struct {
	chainableDifferenceHandler
}

func (handler testableHandler) handle(differences differences) error {
	return nil
}
