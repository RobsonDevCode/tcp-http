package headers

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidHeadersShouldReturnHeaders(t *testing.T) {
	const localHost = "Host: localhost:42069\r\nContent-Type:   application/json\r\n\r\n"
	headers := NewHeaders()
	data := []byte(localHost)

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ := headers.TryGet("host")
	contentType, _ := headers.TryGet("content-type")
	missingKey, _ := headers.TryGet("missingKey")
	require.Equal(t, "localhost:42069", host)
	require.Equal(t, "application/json", contentType)
	require.Equal(t, "", missingKey)

	require.Equal(t, len([]byte(localHost)), n)
	require.True(t, done)
}

func TestInvalidHeaderSpacingShouldReturnErr(t *testing.T) {
	const localHost = " Host : localhost:42069\r\n\r\n"
	headers := NewHeaders()
	data := []byte(localHost)

	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}

func TestInvalidTokenInHeaderShouldReturnErr(t *testing.T) {
	const localHost = "H@st: localhost:42069\r\n\r\n"
	headers := NewHeaders()
	data := []byte(localHost)

	n, done, err := headers.Parse(data)
	require.Error(t, err)
	require.Equal(t, 0, n)
	require.False(t, done)
}

func TestMulitpleValuesOnSameHeaderShouldReturnAppendedValuesOnHeader(t *testing.T) {
	const localHost = "Host: localhost:42069\r\nContent-Type: application/json\r\nContent-Type: application/xml\r\n\r\n"
	headers := NewHeaders()
	data := []byte(localHost)

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ := headers.TryGet("host")
	contentType, _ := headers.TryGet("content-type")
	missingKey, _ := headers.TryGet("missingKey")
	require.Equal(t, "localhost:42069", host)
	require.Equal(t, "application/json,application/xml", contentType)
	require.Equal(t, "", missingKey)

	require.Equal(t, len([]byte(localHost)), n)
	require.True(t, done)
}

func TestMulitpleSameValuesOnSameHeaderShouldReturnAppendedValuesOnHeader(t *testing.T) {
	const localHost = "Host: localhost:42069\r\nContent-Type: application/json\r\nContent-Type: application/json\r\n\r\n"
	headers := NewHeaders()
	data := []byte(localHost)

	n, done, err := headers.Parse(data)

	require.NoError(t, err)
	require.NotNil(t, headers)

	host, _ := headers.TryGet("host")
	contentType, _ := headers.TryGet("content-type")
	missingKey, _ := headers.TryGet("missingKey")

	require.Equal(t, "localhost:42069", host)
	require.Equal(t, "application/json,application/json", contentType)
	require.Equal(t, "", missingKey)

	require.Equal(t, len([]byte(localHost)), n)
	require.True(t, done)
}
