package loghooks

import (
	"fmt"
	"testing"

	"github.com/0xProject/0x-mesh/constants"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type myStruct struct {
	myInt    int
	myString string
}

func TestGetTypeForValue(t *testing.T) {
	testCases := []struct {
		input    interface{}
		expected string
	}{
		{
			input:    true,
			expected: "bool",
		},
		{
			input:    int(42),
			expected: "number",
		},
		{
			input:    int8(42),
			expected: "number",
		},
		{
			input:    int16(42),
			expected: "number",
		},
		{
			input:    int32(42),
			expected: "number",
		},
		{
			input:    int64(42),
			expected: "number",
		},
		{
			input:    uint(42),
			expected: "number",
		},
		{
			input:    uint8(42),
			expected: "number",
		},
		{
			input:    uint16(42),
			expected: "number",
		},
		{
			input:    uint32(42),
			expected: "number",
		},
		{
			input:    uint64(42),
			expected: "number",
		},
		{
			input:    float32(42),
			expected: "number",
		},
		{
			input:    float64(42),
			expected: "number",
		},
		{
			input:    "foo",
			expected: "string",
		},
		{
			input:    complex64(42i + 7),
			expected: "string",
		},
		{
			input:    complex128(42i + 7),
			expected: "string",
		},
		{
			input:    func() {},
			expected: "string",
		},
		{
			input:    make(chan struct{}),
			expected: "string",
		},
		{
			input:    []int{},
			expected: "array",
		},
		{
			input:    [...]int{},
			expected: "array",
		},
		{
			input:    myStruct{},
			expected: "loghooks_myStruct",
		},
		{
			// Implements encoding.TextUnmarshaler but not json.Marshaler.
			input:    constants.NullAddress,
			expected: "string",
		},
		{
			// We don't expect the case of anonymous structs to come up often, but we
			// do handle it correcly. " ", "{", "}", and ";" are allowed in
			// Elasticsearch. The resulting string is ugly but at least it is
			// guaranteed to prevent field mapping conflicts.
			input: struct {
				myInt    int
				myString string
			}{},
			expected: "struct { myInt int; myString string }",
		},
	}

	for _, testCase := range testCases {
		testCaseInfo := fmt.Sprintf("input: (%T) %v", testCase.input, testCase.input)
		actual, err := getTypeForValue(testCase.input)
		require.NoError(t, err, testCaseInfo)
		assert.Equal(t, testCase.expected, actual, testCaseInfo)
	}
}

func TestKeySuffixHookWithNestedMapType(t *testing.T) {
	hook := NewKeySuffixHook()
	entry := &log.Entry{
		Data: log.Fields{
			"myMap": map[string]int{"one": 1},
		},
	}
	require.NoError(t, hook.Fire(entry))
	expectedData := log.Fields{
		"myMap_json_string": `{"one":1}`,
	}
	assert.Equal(t, expectedData, entry.Data)
}
