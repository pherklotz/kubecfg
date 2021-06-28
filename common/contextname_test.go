package common

import (
	"fmt"
	"strings"
	"testing"
)

func TestRandomNameProvider_GetName(t *testing.T) {
	nameProvider := RandomNameProvider{}
	name := nameProvider.GetName()
	if name == "" {
		t.Fatal("Expected a non empty name")
	}
	if !strings.Contains(name, "-") {
		t.Fatal("Expected name with a hyphen but got: ", name)
	}
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
			if name != tt.expected {
				t.Errorf("got '%s' but expected '%s'", name, tt.expected)
			}
		})
	}
}
