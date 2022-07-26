package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewDefaultPage(t *testing.T) {
	data := make([]interface{}, 10)
	data[0] = map[string]interface{}{
		"id":       10001,
		"username": "kami",
	}
	data[1] = map[string]interface{}{
		"id":       10002,
		"username": "monk",
	}

	page := New(121, 10, data, 499)

	assert.Equal(t, 10, page.DataSize())
	assert.Equal(t, 121, page.PageNumber())
	assert.Equal(t, 10, page.PageSize())
	assert.Equal(t, 50, page.TotalPages())
	assert.Equal(t, 1200, page.Offset())
	assert.Equal(t, data, page.Data())
	assert.True(t, page.HasNext())
	assert.True(t, page.HasData())

	t.Logf("Data: %v", page.Data())

	page = New(492, 10, data, 499)
	assert.False(t, page.HasNext())
}
