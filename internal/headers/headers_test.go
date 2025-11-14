package headers

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\nFooFoo:      barbar     \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	host, _ := headers.Get("Host")
	assert.Equal(t, "localhost:42069", host)
	fooFoo, _ := headers.Get("FooFoo")
	assert.Equal(t, "barbar", fooFoo)
	missingKey, _ := headers.Get("MissingKey")
	assert.Equal(t, "", missingKey)
	assert.Equal(t, 51, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Invalid Token
	headers = NewHeaders()
	data = []byte("H©st : localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	// Test: Multiple Values
	headers = NewHeaders()
	data = []byte("Set-Person: lane-loves-go\r\nSet-Person: prime-loves-zig\r\nSet-Person: tj-loves-ocaml\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	setPerson, _ := headers.Get("Set-Person")
	assert.Equal(t, "lane-loves-go, prime-loves-zig, tj-loves-ocaml", setPerson)
	assert.True(t, done)
}
