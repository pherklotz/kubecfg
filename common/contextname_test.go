package common

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomNameProvider_GetName(t *testing.T) {
	nameProvider := RandomNameProvider{}
	name := nameProvider.GetName()
	assert.NotEqual(t, "", name)
	assert.Contains(t, name, "-", "Expected name with a hyphen but got: ", name)
}
func TestSeedNameProvider_GetName(t *testing.T) {
	const seed = "testseed"
	nameProvider := SeedNameProvider{Seed: seed}
	var tests = []struct {
		expected string
	}{
		{seed},
		{seed + "-1"},
		{seed + "-2"},
	}
	for index, tt := range tests {
		testname := fmt.Sprintf("Expect '%s' on %d call of GetName", tt.expected, index+1)
		t.Run(testname, func(t *testing.T) {
			name := nameProvider.GetName()
			assert.Equal(t, tt.expected, name)
		})
	}
}
