package cli

import (
	"fmt"
	"strconv"
)

type Parser interface {
	Match([]string) ([]string, bool)
	Capture([]string) ([]string, map[string]interface{})
}

type EqualParser struct {
	value string
}

func (self EqualParser) Match(values []string) ([]string, bool) {
	if len(values) == 0 {
		return values, false
	}

	if self.value == values[0] {
		return values[1:], true
	}

	return values, false
}

func (self EqualParser) Capture(values []string) ([]string, map[string]interface{}) {
	remainingValues, _ := self.Match(values)
	return remainingValues, nil
}

func (self EqualParser) String() string {
	return fmt.Sprintf("EqualParser '%s'", self.value)
}

type CaptureGroupParser struct {
	value string
}

func (self CaptureGroupParser) Match(values []string) ([]string, bool) {
	if len(values) == 0 {
		return values, false
	}

	return values[1:], true
}

func (self CaptureGroupParser) key() string {
	return self.value[1 : len(self.value)-1]
}

func (self CaptureGroupParser) Capture(values []string) ([]string, map[string]interface{}) {
	if remainingValues, ok := self.Match(values); ok {
		return remainingValues, map[string]interface{}{self.key(): values[0]}
	}

	return values, nil
}

func (self CaptureGroupParser) String() string {
	return fmt.Sprintf("CaptureGroupParser '%s'", self.value)
}

type BoolFlagParser struct {
	pattern      string
	key          string
	omitValue    bool
	defaultValue bool
}

func (self BoolFlagParser) Match(values []string) ([]string, bool) {
	if self.omitValue {
		return flagKeyMatch(self.pattern, values, 0)
	}

	remaining, value, ok := flagKeyValueMatch(self.pattern, values, 0)
	if !ok {
		return remaining, false
	}

	// Check that value is a valid boolean
	if _, err := strconv.ParseBool(value); err != nil {
		return remaining, false
	}

	return remaining, true
}

func (self BoolFlagParser) Capture(values []string) ([]string, map[string]interface{}) {
	if self.omitValue {
		remaining, ok := flagKeyMatch(self.pattern, values, 0)
		return remaining, map[string]interface{}{self.key: ok}
	}

	remaining, value, ok := flagKeyValueMatch(self.pattern, values, 0)
	if !ok {
		return remaining, map[string]interface{}{self.key: self.defaultValue}
	}

	b, _ := strconv.ParseBool(value)
	return remaining, map[string]interface{}{self.key: b}
}

func (self BoolFlagParser) String() string {
	return fmt.Sprintf("BoolFlagParser '%s'", self.pattern)
}

type StringFlagParser struct {
	pattern      string
	key          string
	defaultValue string
}

func (self StringFlagParser) Match(values []string) ([]string, bool) {
	remaining, _, ok := flagKeyValueMatch(self.pattern, values, 0)
	return remaining, ok
}

func (self StringFlagParser) Capture(values []string) ([]string, map[string]interface{}) {
	remaining, value, ok := flagKeyValueMatch(self.pattern, values, 0)
	if !ok {
		return remaining, map[string]interface{}{self.key: self.defaultValue}
	}

	return remaining, map[string]interface{}{self.key: value}
}

func (self StringFlagParser) String() string {
	return fmt.Sprintf("StringFlagParser '%s'", self.pattern)
}

type IntFlagParser struct {
	pattern      string
	key          string
	defaultValue int64
}

func (self IntFlagParser) Match(values []string) ([]string, bool) {
	remaining, value, ok := flagKeyValueMatch(self.pattern, values, 0)
	if !ok {
		return remaining, false
	}

	// Check that value is a valid integer
	if _, err := strconv.ParseInt(value, 10, 64); err != nil {
		return remaining, false
	}

	return remaining, true
}

func (self IntFlagParser) Capture(values []string) ([]string, map[string]interface{}) {
	remaining, value, ok := flagKeyValueMatch(self.pattern, values, 0)
	if !ok {
		return remaining, map[string]interface{}{self.key: self.defaultValue}
	}

	n, _ := strconv.ParseInt(value, 10, 64)
	return remaining, map[string]interface{}{self.key: n}
}

func (self IntFlagParser) String() string {
	return fmt.Sprintf("IntFlagParser '%s'", self.pattern)
}

type StringSliceFlagParser struct {
	pattern      string
	key          string
	defaultValue []string
}

func (self StringSliceFlagParser) Match(values []string) ([]string, bool) {
	if len(values) < 2 {
		return values, false
	}

	var remainingValues []string

	for i := 0; i < len(values); i++ {
		if values[i] == self.pattern && i+1 < len(values) {
			i++
			continue
		}
		remainingValues = append(remainingValues, values[i])
	}

	return remainingValues, len(values) != len(remainingValues)
}

func (self StringSliceFlagParser) Capture(values []string) ([]string, map[string]interface{}) {
	remainingValues, ok := self.Match(values)
	if !ok {
		return values, map[string]interface{}{self.key: self.defaultValue}
	}

	var captured []string

	for i := 0; i < len(values); i++ {
		if values[i] == self.pattern && i+1 < len(values) {
			captured = append(captured, values[i+1])
		}
	}

	return remainingValues, map[string]interface{}{self.key: captured}
}

func (self StringSliceFlagParser) String() string {
	return fmt.Sprintf("StringSliceFlagParser '%s'", self.pattern)
}

type FlagParser struct {
	parsers []Parser
}

func (self FlagParser) Match(values []string) ([]string, bool) {
	remainingValues := values

	for _, parser := range self.parsers {
		remainingValues, _ = parser.Match(remainingValues)
	}
	return remainingValues, true
}

func (self FlagParser) Capture(values []string) ([]string, map[string]interface{}) {
	captured := map[string]interface{}{}
	remainingValues := values

	for _, parser := range self.parsers {
		var data map[string]interface{}
		remainingValues, data = parser.Capture(remainingValues)
		for key, value := range data {
			captured[key] = value
		}
	}

	return remainingValues, captured
}

func (self FlagParser) String() string {
	return fmt.Sprintf("FlagParser %v", self.parsers)
}

type ShortCircuitParser struct {
	parsers []Parser
}

func (self ShortCircuitParser) Match(values []string) ([]string, bool) {
	remainingValues := values

	for _, parser := range self.parsers {
		var ok bool
		remainingValues, ok = parser.Match(remainingValues)
		if ok {
			return remainingValues, true
		}
	}

	return remainingValues, false
}

func (self ShortCircuitParser) Capture(values []string) ([]string, map[string]interface{}) {
	if len(self.parsers) == 0 {
		return values, nil
	}

	for _, parser := range self.parsers {
		if _, ok := parser.Match(values); ok {
			return parser.Capture(values)
		}
	}

	// No parsers matched at this point,
	// just return the capture value of the first one
	return self.parsers[0].Capture(values)
}

func (self ShortCircuitParser) String() string {
	return fmt.Sprintf("ShortCircuitParser %v", self.parsers)
}

type CompleteParser struct {
	parsers []Parser
}

func (self CompleteParser) Match(values []string) ([]string, bool) {
	remainingValues := copySlice(values)

	for _, parser := range self.parsers {
		var ok bool
		remainingValues, ok = parser.Match(remainingValues)
		if !ok {
			return remainingValues, false
		}
	}

	return remainingValues, len(remainingValues) == 0
}

func (self CompleteParser) Capture(values []string) ([]string, map[string]interface{}) {
	remainingValues := copySlice(values)
	data := map[string]interface{}{}

	for _, parser := range self.parsers {
		var captured map[string]interface{}
		remainingValues, captured = parser.Capture(remainingValues)
		for key, value := range captured {
			data[key] = value
		}
	}

	return remainingValues, data
}

func (self CompleteParser) String() string {
	return fmt.Sprintf("CompleteParser %v", self.parsers)
}

func flagKeyValueMatch(key string, values []string, index int) ([]string, string, bool) {
	if index > len(values)-2 {
		return values, "", false
	}

	if values[index] == key {
		value := values[index+1]
		remaining := append(copySlice(values[:index]), values[index+2:]...)
		return remaining, value, true
	}

	return flagKeyValueMatch(key, values, index+1)
}

func flagKeyMatch(key string, values []string, index int) ([]string, bool) {
	if index > len(values)-1 {
		return values, false
	}

	if values[index] == key {
		remaining := append(copySlice(values[:index]), values[index+1:]...)
		return remaining, true
	}

	return flagKeyMatch(key, values, index+1)
}

func copySlice(a []string) []string {
	b := make([]string, len(a))
	copy(b, a)
	return b
}
