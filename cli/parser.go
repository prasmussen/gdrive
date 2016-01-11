package cli

import (
    "fmt"
    "strconv"
)

type Parser interface {
    Match([]string) ([]string, bool)
    Capture([]string) ([]string, map[string]string)
}

type CompleteParser struct {
    parsers []Parser
}

func (self CompleteParser) Match(values []string) ([]string, bool) {
    remainingValues := values

    for _, parser := range self.parsers {
        var ok bool
        remainingValues, ok = parser.Match(remainingValues)
        if !ok {
            return remainingValues, false
        }
    }

    return remainingValues, len(remainingValues) == 0
}

func (self CompleteParser) Capture(values []string) ([]string, map[string]string) {
    remainingValues := values
    data := map[string]string{}

    for _, parser := range self.parsers {
        var captured map[string]string
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

func (self EqualParser) Capture(values []string) ([]string, map[string]string) {
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
    return self.value[1:len(self.value) - 1]
}

func (self CaptureGroupParser) Capture(values []string) ([]string, map[string]string) {
    if remainingValues, ok := self.Match(values); ok {
        return remainingValues, map[string]string{self.key(): values[0]}
    }

    return values, nil
}

func (self CaptureGroupParser) String() string {
    return fmt.Sprintf("CaptureGroupParser '%s'", self.value)
}



type BoolFlagParser struct {
    pattern string
    key string
    omitValue bool
    defaultValue bool
}

func (self BoolFlagParser) Match(values []string) ([]string, bool) {
    if self.omitValue {
        if len(values) == 0 {
            return values, false
        }

        if self.pattern == values[0] {
            return values[1:], true
        }

        return values, false
    } else {
        if len(values) < 2 {
            return values, false
        }

        if self.pattern != values[0] {
            return values, false
        }

        // Check that value is a valid boolean
        if _, err := strconv.ParseBool(values[1]); err != nil {
            return values, false
        }

        return values[2:], true
    }
}

func (self BoolFlagParser) Capture(values []string) ([]string, map[string]string) {
    remainingValues, ok := self.Match(values)
    if !ok && !self.omitValue {
        return remainingValues, map[string]string{self.key: fmt.Sprintf("%t", self.defaultValue)}
    }
    return remainingValues, map[string]string{self.key: fmt.Sprintf("%t", ok)}
}

func (self BoolFlagParser) String() string {
    return fmt.Sprintf("BoolFlagParser '%s'", self.pattern)
}

type StringFlagParser struct {
    pattern string
    key string
    defaultValue string
}

func (self StringFlagParser) Match(values []string) ([]string, bool) {
    if len(values) < 2 {
        return values, false
    }

    if self.pattern != values[0] {
        return values, false
    }

    return values[2:], true
}

func (self StringFlagParser) Capture(values []string) ([]string, map[string]string) {
    remainingValues, ok := self.Match(values)
    if ok {
        return remainingValues, map[string]string{self.key: values[1]}
    }

    return values, map[string]string{self.key: self.defaultValue}
}

func (self StringFlagParser) String() string {
    return fmt.Sprintf("StringFlagParser '%s'", self.pattern)
}

type IntFlagParser struct {
    pattern string
    key string
    defaultValue int64
}

func (self IntFlagParser) Match(values []string) ([]string, bool) {
    if len(values) < 2 {
        return values, false
    }

    if self.pattern != values[0] {
        return values, false
    }

    // Check that value is a valid integer
    if _, err := strconv.ParseInt(values[1], 10, 64); err != nil {
        return values, false
    }

    return values[2:], true
}

func (self IntFlagParser) Capture(values []string) ([]string, map[string]string) {
    remainingValues, ok := self.Match(values)
    if ok {
        return remainingValues, map[string]string{self.key: values[1]}
    }

    return values, map[string]string{self.key: fmt.Sprintf("%d", self.defaultValue)}
}

func (self IntFlagParser) String() string {
    return fmt.Sprintf("IntFlagParser '%s'", self.pattern)
}


type FlagParser struct {
    parsers []Parser
}

func (self FlagParser) Match(values []string) ([]string, bool) {
    remainingValues := values
    var oneOrMoreMatches bool

    for _, parser := range self.parsers {
        var ok bool
        remainingValues, ok = parser.Match(remainingValues)
        if ok {
            oneOrMoreMatches = true
        }
    }

    // Recurse while we have one or more matches
    if oneOrMoreMatches {
        return self.Match(remainingValues)
    }

    return remainingValues, true
}

func (self FlagParser) Capture(values []string) ([]string, map[string]string) {
    data := map[string]string{}
    remainingValues := values

    for _, parser := range self.parsers {
        var captured map[string]string
        remainingValues, captured = parser.Capture(remainingValues)
        for key, value := range captured {
            // Skip value if it already exists and new value is an empty string
            if _, exists := data[key]; exists && value == "" {
                continue
            }

            data[key] = value
        }
    }
    return remainingValues, data
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

func (self ShortCircuitParser) Capture(values []string) ([]string, map[string]string) {
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
