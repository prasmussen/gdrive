package cli

import (
    "fmt"
    "regexp"
    "strings"
)

type Flags map[string][]Flag

var handlers []*Handler

type Handler struct {
    Pattern string
    Flags Flags
    Callback func(Context)
    Description string
}

func (self *Handler) getParser() Parser {
    var parsers []Parser

    for _, pattern := range splitPattern(self.Pattern) {
        if isOptional(pattern) {
            name := optionalName(pattern)
            parser := getFlagParser(self.Flags[name])
            parsers = append(parsers, parser)
        } else if isCaptureGroup(pattern) {
            parsers = append(parsers, CaptureGroupParser{pattern})
        } else {
            parsers = append(parsers, EqualParser{pattern})
        }
    }

    return CompleteParser{parsers}
}

func SetHandlers(h []*Handler) {
    handlers = h
}

func AddHandler(pattern string, flags Flags, callback func(Context), desc string) {
    handlers = append(handlers, &Handler{
        Pattern: pattern,
        Flags: flags,
        Callback: callback,
        Description: desc,
    })
}

func findHandler(args []string) *Handler {
    for _, h := range handlers {
        if _, ok := h.getParser().Match(args); ok {
            return h
        }
    }
    return nil
}


func Handle(args []string) bool {
    h := findHandler(args)
    if h == nil {
        return false
    }

    _, data := h.getParser().Capture(args)
    fmt.Println(data)
    ctx := Context{
        args: data,
        handlers: handlers,
    }
    h.Callback(ctx)
    return true
}

func filterHandlers(handlers []*Handler, prefix string) []*Handler {
    matches := []*Handler{}

    for _, h := range handlers {
        pattern := strings.Join(stripOptionals(splitPattern(h.Pattern)), " ")
        if strings.HasPrefix(pattern, prefix) {
            matches = append(matches, h)
        }
    }

    return matches
}


// Split on spaces but ignore spaces inside <...> and [...]
func splitPattern(pattern string) []string {
    re := regexp.MustCompile(`(<[^>]+>|\[[^\]]+]|\S+)`)
    matches := []string{}

    for _, value := range re.FindAllStringSubmatch(pattern, -1) {
        matches = append(matches, value[1]) 
    }

    return matches
}

func isCaptureGroup(arg string) bool {
    return strings.HasPrefix(arg, "<") && strings.HasSuffix(arg, ">")
}

func isOptional(arg string) bool {
    return strings.HasPrefix(arg, "[") && strings.HasSuffix(arg, "]")
}

func optionalName(s string) string {
    return s[1:len(s) - 1]
}

// Strip optional groups from pattern
func stripOptionals(pattern []string) []string {
    newArgs := []string{}

    for _, arg := range pattern {
        if !isOptional(arg) {
            newArgs = append(newArgs, arg)
        }
    }
    return newArgs
}
