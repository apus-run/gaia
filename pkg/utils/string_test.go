package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	assert.Equal(t, "test", EscapeQuote("test"))
	assert.Equal(t, "test", EscapeQuote("`test`"))
	assert.Equal(t, `\'test\'`, EscapeQuote("'test'"))
	assert.Equal(t, `\"test\"`, EscapeQuote(`"test"`))
	assert.Equal(t, `\\test\\`, EscapeQuote(`\test\`))
}
