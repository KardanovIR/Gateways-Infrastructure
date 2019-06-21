package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNodeClient_replaceQuotesFromSides(t *testing.T) {
	s := "\"c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73\""
	assert.Equal(t, "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73", replaceQuotesFromSides(s))
	s2 := "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73"
	assert.Equal(t, "c99ce66bd0cebf45f97aba2f48912583562a7cc8fdf4c89079608517b1955c73", replaceQuotesFromSides(s2))
}
