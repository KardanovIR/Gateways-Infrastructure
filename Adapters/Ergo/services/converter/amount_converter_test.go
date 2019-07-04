package converter

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToTargetAmount(t *testing.T) {
	Init(context.Background(), 8)
	assert.Equal(t, uint64(1234560), ToNodeAmount(123456))
	assert.Equal(t, uint64(123456), ToTargetAmount(1234560))
	assert.Equal(t, uint64(12345), ToTargetAmount(123456))
}
