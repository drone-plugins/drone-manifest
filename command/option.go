package command

type Option func(*Command)

func WithUsername(val string) Option {
	return func(t *Command) {
		t.username = val
	}
}

func WithPassword(val string) Option {
	return func(t *Command) {
		t.password = val
	}
}

func WithSpec(val string) Option {
	return func(t *Command) {
		t.spec = val
	}
}

func WithPlatforms(val []string) Option {
	return func(t *Command) {
		t.platforms = val
	}
}

func WithTarget(val string) Option {
	return func(t *Command) {
		t.target = val
	}
}

func WithTemplate(val string) Option {
	return func(t *Command) {
		t.template = val
	}
}

func WithPath(val string) Option {
	return func(t *Command) {
		t.path = val
	}
}

func IgnoreMissing() Option {
	return func(t *Command) {
		t.ignoreMissing = true
	}
}
