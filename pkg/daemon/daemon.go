package daemon

type Config struct {
	Namespace  string
	Webhook    string
	NotifyType string
}

type Daemon struct {
	config Config
}

func NewConfig(namespace, webhook, notifyType string) (Config, error) {
	return Config{Namespace: namespace, Webhook: webhook, NotifyType: notifyType}, nil
}

func NewDaemon(c Config) (Daemon, error) {
	return Daemon{config: c}, nil
}

func (d Daemon) Start() {

}
