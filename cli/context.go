package cli

type Context struct {
	args     Arguments
	handlers []*Handler
}

func (self Context) Args() Arguments {
	return self.args
}

func (self Context) Handlers() []*Handler {
	return self.handlers
}

type Arguments map[string]interface{}

func (self Arguments) String(key string) string {
	return self[key].(string)
}

func (self Arguments) Int64(key string) int64 {
	return self[key].(int64)
}

func (self Arguments) Bool(key string) bool {
	return self[key].(bool)
}

func (self Arguments) StringSlice(key string) []string {
	return self[key].([]string)
}
