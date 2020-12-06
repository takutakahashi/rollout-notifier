package notify

type Notifier interface {
	Start(target string) error
	Finish(target string) error
}

type Config struct {
	Path string
}

func NewNotify(t, configPath string) (Notifier, error) {
	switch t {
	case "slack":
		return NewSlackNotify(configPath)
	default:
		return NoopNotify{}, nil
	}
}
