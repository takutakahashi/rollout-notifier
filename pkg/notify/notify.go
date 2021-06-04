package notify

type Notifier interface {
	Start(target, comment string) error
	Finish(target, comment string) error
	Failed(target, comment string) error
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
