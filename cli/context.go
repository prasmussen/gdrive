package cli

import (
    "strconv"
)

type Context struct {
    args Arguments
    handlers []*Handler
}

func (self Context) Args() Arguments {
    return self.args
}

func (self Context) Handlers() []*Handler {
    return self.handlers
}

func (self Context) FilterHandlers(prefix string) []*Handler {
    return filterHandlers(self.handlers, prefix)
}

type Arguments map[string]string

func (self Arguments) String(key string) string {
    value, _ := self[key]
    return value
}

func (self Arguments) Int64(key string) int64 {
    value, _ := self[key]
    n, _ := strconv.ParseInt(value, 10, 64)
    return n
}

func (self Arguments) Bool(key string) bool {
    value, _ := self[key]
    b, _ := strconv.ParseBool(value)
    return b
}
