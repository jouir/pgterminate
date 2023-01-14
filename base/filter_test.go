package base

import (
	"fmt"
	"regexp"
	"testing"
)

func TestIncludeFilter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		patterns []string
		wanted   bool
	}{
		{"No filter", "test", nil, true},
		{"Empty filter", "test", []string{""}, true},
		{"Single pattern matching", "test", []string{"test"}, true},
		{"Multiple patterns matching", "test", []string{"test", "postgres"}, true},
		{"Single pattern with no match", "nomatch", []string{"test"}, false},
		{"Multiple patterns with no match", "nomatch", []string{"test", "postgres"}, false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			f := NewIncludeFilter(tc.patterns)

			if got := f.Include(tc.value); got != tc.wanted {
				t.Errorf("Included must be %t for patterns '%s'", tc.wanted, tc.patterns)
			} else {
				t.Logf("Included is %t for patterns '%s'", tc.wanted, tc.patterns)
			}
		})
	}
}

func TestIncludeFilterRegex(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		regex  string
		wanted bool
	}{
		{"No filter", "test", "", true},
		{"String pattern matching", "test", "test", true},
		{"Regex patterns matching", "test", "^t(.*)$", true},
		{"String pattern with no match", "nomatch", "test", false},
		{"Regex patterns with no match", "nomatch", "^t(.*)$", false},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			compiledRegex, err := regexp.Compile(tc.regex)
			if err != nil {
				t.Fatalf("Regex '%s' doesn't compile: %v", tc.regex, err)
			}

			f := NewIncludeFilterRegex(compiledRegex)
			if got := f.Include(tc.value); got != tc.wanted {
				t.Errorf("Included must be %t for regex '%s'", tc.wanted, tc.regex)
			} else {
				t.Logf("Included is %t for regex '%s'", tc.wanted, tc.regex)
			}
		})
	}
}

func TestExcludeFilter(t *testing.T) {
	tests := []struct {
		name     string
		value    string
		patterns []string
		wanted   bool
	}{
		{"No filter", "test", nil, true},
		{"Empty filter", "test", []string{""}, true},
		{"Single pattern matching", "test", []string{"test"}, false},
		{"Multiple patterns matching", "test", []string{"test", "postgres"}, false},
		{"Single pattern with no match", "nomatch", []string{"test"}, true},
		{"Multiple patterns with no match", "nomatch", []string{"test", "postgres"}, true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			f := NewExcludeFilter(tc.patterns)
			if got := f.Include(tc.value); got != tc.wanted {
				t.Errorf("Included must be %t for patterns '%s'", tc.wanted, tc.patterns)
			} else {
				t.Logf("Included is %t for patterns '%s'", tc.wanted, tc.patterns)
			}
		})
	}
}

func TestExcludeFilterRegex(t *testing.T) {
	tests := []struct {
		name   string
		value  string
		regex  string
		wanted bool
	}{
		{"No filter", "test", "", true},
		{"String pattern matching", "test", "test", false},
		{"Regex patterns matching", "test", "^t(.*)$", false},
		{"String pattern with no match", "nomatch", "test", true},
		{"Regex patterns with no match", "nomatch", "^t(.*)$", true},
	}

	for _, tc := range tests {
		t.Run(fmt.Sprintf(tc.name), func(t *testing.T) {
			compiledRegex, err := regexp.Compile(tc.regex)
			if err != nil {
				t.Fatalf("Regex '%s' doesn't compile: %v", tc.regex, err)
			}

			f := NewExcludeFilterRegex(compiledRegex)
			if got := f.Include(tc.value); got != tc.wanted {
				t.Errorf("Included must be %t for regex '%s'", tc.wanted, tc.regex)
			} else {
				t.Logf("Included is %t for regex '%s'", tc.wanted, tc.regex)
			}
		})
	}
}
