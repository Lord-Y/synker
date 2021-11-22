// Package tools assemble useful functions used by other packages
package tools

import (
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"
)

func TestRandString(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value  int
		actual int
	}{
		{
			value:  0,
			actual: 1,
		},
		{
			value:  -20,
			actual: 1,
		},
		{
			value:  1,
			actual: 1,
		},
		{
			value:  4,
			actual: 4,
		},
		{
			value:  500,
			actual: 500,
		},
	}

	for _, tc := range tests {
		z := RandString(tc.value)
		assert.Equal(len(z), tc.actual)
	}
}

func TestRandStringInt(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value  int
		actual int
	}{
		{
			value:  0,
			actual: 1,
		},
		{
			value:  -20,
			actual: 1,
		},
		{
			value:  1,
			actual: 1,
		},
		{
			value:  4,
			actual: 4,
		},
		{
			value:  500,
			actual: 500,
		},
	}

	for _, tc := range tests {
		z := RandStringInt(tc.value)
		assert.Equal(len(z), tc.actual)
	}
}

func TestCaptureOutput(t *testing.T) {
	assert := assert.New(t)

	f1 := func() {
		log.Print("f1")
	}

	output := CaptureOutput(f1)
	assert.Contains(output, "f1")
}

func TestGetPagination(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		page             int
		start            int
		end              int
		rangeLimit       int
		actualStartLimit int
		actualEndLimit   int
	}{
		{
			page:             0,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             1,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             2,
			start:            2,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 20,
			actualEndLimit:   20,
		},
		{
			page:             -1,
			start:            1,
			end:              20,
			rangeLimit:       20,
			actualStartLimit: 0,
			actualEndLimit:   20,
		},
		{
			page:             -1,
			start:            -1,
			end:              -20,
			rangeLimit:       -20,
			actualStartLimit: 0,
			actualEndLimit:   1,
		},
	}

	for _, tc := range tests {
		sl, el := GetPagination(tc.page, tc.start, tc.end, tc.rangeLimit)
		assert.Equal(sl, tc.actualStartLimit)
		assert.Equal(el, tc.actualEndLimit)
	}
}

func TestCheckIsFile(t *testing.T) {
	assert := assert.New(t)

	f, err := os.CreateTemp(os.TempDir(), fake.CharactersN(5))
	if err != nil {
		assert.Fail("Fail to create temp file")
		return
	}
	defer os.Remove(f.Name())

	z := CheckIsFile(f.Name())
	assert.Nil(z)

	z = CheckIsFile(fmt.Sprintf("%s/%s", os.TempDir(), fake.CharactersN(10)))
	assert.Error(z)

	directory, err := os.MkdirTemp(os.TempDir(), fake.CharactersN(10))
	if err != nil {
		assert.Fail("Fail to create temp directory")
		return
	}
	defer os.RemoveAll(directory)

	z = CheckIsFile(directory)
	assert.Error(z)
}

func TestRandomValueFromSlice(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		slice []string
	}{
		{
			slice: []string{"a", "1", "", "v"},
		},
		{
			slice: []string{"a", "1", "v"},
		},
	}

	for i, tc := range tests {
		z := RandomValueFromSlice(tc.slice)
		if i == 0 {
			switch len(z) {
			case 0:
				assert.Equal(0, len(z))
			case 1:
				assert.Equal(1, len(z))
			default:
				t.Fail()
			}
		} else {
			assert.Equal(1, len(z))
		}
	}
}

func TestIsJSON(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    "v",
			expected: false,
		},
		{
			value:    `{"a": "b"}`,
			expected: true,
		},
		{
			value:    `{'a': 'b'}`,
			expected: false,
		},
		{
			value:    `"a": "b"`,
			expected: false,
		},
	}

	for _, tc := range tests {
		z := IsJSON(tc.value)
		assert.Equal(tc.expected, z)
	}
}

func TestIsJSONFromBytes(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    "v",
			expected: false,
		},
		{
			value:    `{"a": "b"}`,
			expected: true,
		},
		{
			value:    `{'a': 'b'}`,
			expected: false,
		},
		{
			value:    `"a": "b"`,
			expected: false,
		},
	}

	for _, tc := range tests {
		z := IsJSONFromBytes([]byte(tc.value))
		assert.Equal(tc.expected, z)
	}
}

func TestIsYaml(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    "v",
			expected: false,
		},
		{
			value:    `a: b`,
			expected: true,
		},
		{
			value:    `'a': `,
			expected: true,
		},
		{
			value:    `'a': 'b'`,
			expected: true,
		},
		{
			value:    `"a": "b"`,
			expected: true,
		},
	}

	for _, tc := range tests {
		z := IsYaml(tc.value)
		assert.Equal(tc.expected, z)
	}
}

func TestIsYamlFromBytes(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value    string
		expected bool
	}{
		{
			value:    "v",
			expected: false,
		},
		{
			value:    `a: b`,
			expected: true,
		},
		{
			value:    `'a': `,
			expected: true,
		},
		{
			value:    `'a': 'b'`,
			expected: true,
		},
		{
			value:    `"a": "b"`,
			expected: true,
		},
	}

	for _, tc := range tests {
		z := IsYamlFromBytes([]byte(tc.value))
		assert.Equal(tc.expected, z)
	}
}

func TestInSlice(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		value    string
		array    []string
		expected bool
	}{
		{
			value:    "value",
			array:    []string{"a", "b", "value"},
			expected: true,
		},
		{
			value:    "v",
			array:    []string{"a", "b", "value"},
			expected: false,
		},
	}

	for _, tc := range tests {
		z := InSlice(tc.value, tc.array)
		assert.Equal(tc.expected, z)
	}
}
