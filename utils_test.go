package godm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseLine(t *testing.T) {
	t.Run("Parses basic without units", func(t *testing.T) {
		k, v, err := parseLine("KEY=VALUE")
		assert.Nil(t, err)
		assert.Equal(t, "KEY", k)
		assert.Equal(t, "VALUE", v)
	})

	t.Run("Parses basic with spaces", func(t *testing.T) {
		k, v, err := parseLine(" KEY = VALUE ")
		assert.Nil(t, err)
		assert.Equal(t, "KEY", k)
		assert.Equal(t, "VALUE", v)

		k, v, err = parseLine(" KEY =          VALUE ")
		assert.Nil(t, err)
		assert.Equal(t, "KEY", k)
		assert.Equal(t, "VALUE", v)
	})

	t.Run("Parses MESSAGE_ID example", func(t *testing.T) {
		k, v, err := parseLine("MESSAGE_ID = OPM 201113719185")
		assert.Nil(t, err)
		assert.Equal(t, "MESSAGE_ID", k)
		assert.Equal(t, "OPM 201113719185", v)
	})

	t.Run("Parses basic with units", func(t *testing.T) {
		k, v, err := parseLine("KEY=VALUE [ms]")
		assert.Nil(t, err)
		assert.Equal(t, "KEY", k)
		assert.Equal(t, "VALUE", v)

		k, v, err = parseLine("GM = 398600.4415 [km**3/s**2]")
		assert.Nil(t, err)
		assert.Equal(t, "GM", k)
		assert.Equal(t, "398600.4415", v)
	})

	t.Run("Parses COMMENT", func(t *testing.T) {
		k, v, err := parseLine("COMMENT      1996-01-01 00:00:00.000")
		assert.Nil(t, err)
		assert.Equal(t, "COMMENT", k)
		assert.Equal(t, "1996-01-01 00:00:00.000", v)
	})
}

func TestDetectLineEnding(t *testing.T) {
	t.Run("Detects CRLF", func(t *testing.T) {
		assert.Equal(t, CRLF, detectLineEnding("a\r\nb\r\nc\r\nd"))
	})

	t.Run("Detects LFCR", func(t *testing.T) {
		assert.Equal(t, LFCR, detectLineEnding("a\n\rb\n\rc\n\rd"))
	})

	t.Run("Detects CR", func(t *testing.T) {
		assert.Equal(t, CR, detectLineEnding("a\rb\rc\rd"))
	})

	t.Run("Detects LF", func(t *testing.T) {
		assert.Equal(t, LF, detectLineEnding("a\nb\nc\nd"))
	})
}
