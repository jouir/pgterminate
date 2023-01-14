package base

import (
	"fmt"
	"reflect"
	"regexp"
)

// Filter interface to tell if a string should be included or not
type Filter interface {
	Include(string) bool
	String() string
}

// IncludeFilter to include a string when it's included in a list of strings
type IncludeFilter struct {
	patterns []string
}

// NewIncludeFilter to create an IncludeFilter
func NewIncludeFilter(patterns []string) IncludeFilter {
	return IncludeFilter{
		patterns: patterns,
	}
}

// Include returns true when a string is included in a list of patterns
// Implements the Filter interface
func (f IncludeFilter) Include(s string) bool {
	// No or empty patterns must include
	if f.patterns == nil || reflect.DeepEqual(f.patterns, []string{""}) {
		return true
	}
	return InSlice(s, f.patterns)
}

// String to pretty print an IncludeFilter
// Implements the Filter interface
func (f IncludeFilter) String() string {
	return fmt.Sprintf("<IncludeFilter(%s)>", f.patterns)
}

// IncludeFilterRegex to include a string when it matches a regex
type IncludeFilterRegex struct {
	regex *regexp.Regexp
}

// NewIncludeFilterRegex to create an IncludeFilterRegex
func NewIncludeFilterRegex(regex *regexp.Regexp) IncludeFilterRegex {
	return IncludeFilterRegex{
		regex: regex,
	}
}

// Include returns true when the string matches the regex
// Implements the Filter interface
func (f IncludeFilterRegex) Include(s string) bool {
	if f.regex == nil || f.regex.MatchString(s) {
		return true
	}
	return false
}

// String to pretty print an IncludeFilterRegex
// Implements the Filter interface
func (f IncludeFilterRegex) String() string {
	return fmt.Sprintf("<IncludeFilterRegex(%s)>", f.regex.String())
}

// ExcludeFilter to include a string when it's not included in a list of strings
type ExcludeFilter struct {
	patterns []string
}

// NewExcludeFilter to create an ExcludeFilter
func NewExcludeFilter(patterns []string) ExcludeFilter {
	return ExcludeFilter{
		patterns: patterns,
	}
}

// Include returns true when the string is not included in the patterns
// Implements the Filter interface
func (f ExcludeFilter) Include(s string) bool {
	return !InSlice(s, f.patterns)
}

// String to pretty print an ExcludeFilter
// Implements the Filter interface
func (f ExcludeFilter) String() string {
	return fmt.Sprintf("<ExcludeFilter(%s)>", f.patterns)
}

// ExcludeFilterRegex to include a string when it doesnn't match a regex
type ExcludeFilterRegex struct {
	regex *regexp.Regexp
}

// NewExcludeFilterRegex to create an ExcludeFilterRegex
func NewExcludeFilterRegex(regex *regexp.Regexp) ExcludeFilterRegex {
	return ExcludeFilterRegex{
		regex: regex,
	}
}

// Include returns true when the string doesn't match the regex
// Implements the Filter interface
func (f ExcludeFilterRegex) Include(s string) bool {
	if f.regex == nil || f.regex.MatchString("") {
		return true
	}
	if f.regex.MatchString(s) {
		return false
	}
	return true
}

// String to pretty print an ExcludeFilterRegex
// Implements the Filter interface
func (f ExcludeFilterRegex) String() string {
	return fmt.Sprintf("<ExcludeFilterRegex(%s)>", f.regex.String())
}
