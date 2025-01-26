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

func TestParseIntoKVs(t *testing.T) {
	t.Run("Parses basic", func(t *testing.T) {
		expected := []KV{
			{"KEY", "VALUE"},
		}
		got, err := parseIntoKVs(`KEY=VALUE`)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Parses basic with spaces", func(t *testing.T) {
		expected := []KV{
			{"KEY", "VALUE"},
		}
		got, err := parseIntoKVs(` KEY = VALUE `)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Parses multiple", func(t *testing.T) {
		expected := []KV{
			{"KEY", "VALUE"},
			{"KEY2", "VALUE2"},
		}
		got, err := parseIntoKVs(`KEY=VALUE
KEY2=VALUE2`)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Parses with multiple of same line", func(t *testing.T) {
		expected := []KV{
			{"KEY", "VALUE"},
			{"KEY", "VALUE2"},
		}

		got, err := parseIntoKVs(`KEY=VALUE
KEY=VALUE2`)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})

	t.Run("Parses COMMENT", func(t *testing.T) {
		expected := []KV{
			{"COMMENT", "1996-01-01 00:00:00.000"},
		}
		got, err := parseIntoKVs(`COMMENT      1996-01-01 00:00:00.000`)
		assert.Nil(t, err)
		assert.Equal(t, expected, got)
	})
}

func TestGetODMFields(t *testing.T) {
	test := OPM{}
	fields, err := getODMFields(&test)
	assert.Nil(t, err)
	assert.Equal(t, "CCSDS_OPM_VERS", fields[0].Name)

	fields[0].ReflectVal.SetString("1.0")
	assert.Equal(t, "1.0", test.Header.CcsdsOpmVers)

	foundComment := false
	for _, f := range fields {
		if f.Name == "COMMENT" {
			foundComment = true
			break
		}
	}
	assert.True(t, foundComment, "COMMENT field not found")
}
