package lambda

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSeq(t *testing.T) {
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, Seq(int(1), 6, 1))
	assert.Equal(t, []int{2, 4, 6, 8, 10}, Seq(int(2), 10, 2))
}
